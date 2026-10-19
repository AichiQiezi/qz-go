"""Validate one Interview Drill daily set (interview/data/YYYY-MM-DD.json).

Usage: python3 -I interview/validate_day.py interview/data/2026-10-19.json
Exits non-zero and prints every problem found.

Rules: 5 questions in the fixed topic order, difficulties 2 Medium / 2 Hard / 1 Senior,
ids YYYYMMDD-1..5, English and Chinese blocks that line up one to one (the page shows
each Chinese paragraph under its English one), only the supported markdown subset.
"""
import json
import os
import re
import sys

TOPICS = ["Go", "System Design", "Database", "Middleware", "Java & CS"]
KEYS = ["id", "topic", "difficulty", "title", "title_zh", "question", "question_zh",
        "answer", "answer_zh", "followups", "followups_zh", "vocab"]
ITEM = re.compile(r"^\s*[-*] ")


def blocks(src):
    out, lines, i = [], (src or "").split("\n"), 0
    while i < len(lines):
        line = lines[i]
        if line.startswith("```"):
            i += 1
            while i < len(lines) and not lines[i].startswith("```"):
                i += 1
            i += 1
            out.append(("pre", 0))
        elif ITEM.match(line):
            n = 0
            while i < len(lines) and ITEM.match(lines[i]):
                n += 1
                i += 1
            out.append(("ul", n))
        elif not line.strip():
            i += 1
        else:
            while i < len(lines) and lines[i].strip() and not lines[i].startswith("```") and not ITEM.match(lines[i]):
                i += 1
            out.append(("p", 0))
    return out


def unsupported(text):
    t = re.sub(r"```.*?```", "", text, flags=re.S)
    t = re.sub(r"`[^`]*`", "", t).replace("**", "")
    probs = []
    if "*" in t:
        probs.append("single * (italics or stray asterisk) outside code")
    if re.search(r"^#", t, re.M):
        probs.append("markdown heading")
    if "](" in t:
        probs.append("markdown link")
    if re.search(r"^\s*\d+\. ", t, re.M):
        probs.append("numbered list (use '- ' bullets instead)")
    if re.search(r"^\s*\|", t, re.M):
        probs.append("table")
    return probs


def main():
    path = sys.argv[1]
    date = os.path.basename(path)[:-5]
    errs = []
    d = json.load(open(path, encoding="utf-8"))
    if list(d.keys()) != ["date", "generatedAt", "questions"]:
        errs.append(f"top-level keys must be date, generatedAt, questions; got {list(d.keys())}")
    if d.get("date") != date:
        errs.append(f"date {d.get('date')!r} does not match the file name {date!r}")
    if not re.match(r"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\+08:00$", d.get("generatedAt", "")):
        errs.append("generatedAt must look like 2026-10-19T06:45:00+08:00")
    qs = d.get("questions") or []
    if [q.get("topic") for q in qs] != TOPICS:
        errs.append(f"topics must be exactly {TOPICS} in that order")
    levels = sorted(q.get("difficulty") for q in qs)
    if levels != ["Hard", "Hard", "Medium", "Medium", "Senior"]:
        errs.append(f"difficulties must be 2 Medium, 2 Hard, 1 Senior; got {[q.get('difficulty') for q in qs]}")
    for n, q in enumerate(qs, 1):
        tag = q.get("id", f"#{n}")
        if list(q.keys()) != KEYS:
            errs.append(f"{tag}: keys must be exactly {KEYS} in that order")
            continue
        if q["id"] != f"{date.replace('-', '')}-{n}":
            errs.append(f"{tag}: id should be {date.replace('-', '')}-{n}")
        for k in KEYS:
            if q[k] in ("", None, []):
                errs.append(f"{tag}: empty {k}")
        for en, zh in (("question", "question_zh"), ("answer", "answer_zh")):
            if blocks(q[en]) != blocks(q[zh]):
                errs.append(f"{tag}: {en} vs {zh} blocks differ:\n    en={blocks(q[en])}\n    zh={blocks(q[zh])}")
            if q[en].count("**") != q[zh].count("**"):
                errs.append(f"{tag}: {en}/{zh} bold markers differ")
            for k in (en, zh):
                for p in unsupported(q[k]):
                    errs.append(f"{tag}: {k}: {p}")
        if not (2 <= len(q["followups"]) <= 3) or len(q["followups"]) != len(q["followups_zh"]):
            errs.append(f"{tag}: followups must be 2-3 and match followups_zh")
        if not (3 <= len(q["vocab"]) <= 5) or any(set(v) != {"en", "zh"} or not v["en"] or not v["zh"] for v in q["vocab"]):
            errs.append(f"{tag}: vocab must be 3-5 objects with non-empty en and zh")
        if len(q["answer"]) < 1200:
            errs.append(f"{tag}: answer looks too short ({len(q['answer'])} chars)")
    if errs:
        print(f"FAIL {path}")
        for e in errs:
            print("  -", e)
        sys.exit(1)
    print(f"OK {path}: {[q['difficulty'] for q in qs]}")


main()
