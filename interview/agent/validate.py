"""Validate one Agent course question set before it is published.

Core modules (kind "core") hold 3 questions with ids agent-NN-1..3; daily sets (kind "daily",
id "daily-YYYY-MM-DD") hold 5 questions with ids agent-daily-YYYYMMDD-1..5.
Checks the schema, the markdown subset the page renders, and that the Chinese text lines up with
the English block for block (paragraphs, list items, code blocks).

Usage (from the repo root): python3 -I interview/agent/validate.py interview/agent/data/<id>.json
Exits non-zero and prints every problem found.
"""
import json
import re
import sys

KEYS = ["id", "topic", "difficulty", "title", "title_zh", "question", "question_zh",
        "answer", "answer_zh", "followups", "followups_zh", "vocab"]
ITEM = re.compile(r"^\s*[-*] ")


def blocks(src):
    out, lines, i = [], (src or "").split("\n"), 0
    while i < len(lines):
        l = lines[i]
        if l.startswith("```"):
            buf = []
            i += 1
            while i < len(lines) and not lines[i].startswith("```"):
                buf.append(lines[i])
                i += 1
            i += 1
            out.append(("pre", len(buf)))
        elif ITEM.match(l):
            n = 0
            while i < len(lines) and ITEM.match(lines[i]):
                n += 1
                i += 1
            out.append(("ul", n))
        elif not l.strip():
            i += 1
        else:
            while i < len(lines) and lines[i].strip() and not lines[i].startswith("```") and not ITEM.match(lines[i]):
                i += 1
            out.append(("p", 0))
    return out


def shape(bs):
    return [(t, n if t == "ul" else 0) for t, n in bs]


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
    if "|" in t and re.search(r"^\s*\|", t, re.M):
        probs.append("table")
    return probs


def main():
    path = sys.argv[1]
    want = sys.argv[2] if len(sys.argv) > 2 else path.rsplit("/", 1)[-1].removesuffix(".json")
    errs = []
    d = json.load(open(path, encoding="utf-8"))
    for k in ("id", "order", "kind", "title", "title_zh", "summary", "summary_zh", "generatedAt", "questions"):
        if k not in d or d[k] in ("", None):
            errs.append(f"module: missing {k}")
    if want and d.get("id") != want:
        errs.append(f"module id {d.get('id')!r} != {want!r}")
    qs = d.get("questions") or []
    daily = d.get("kind") == "daily"
    want_n = 5 if daily else 3
    if len(qs) != want_n:
        errs.append(f"expected {want_n} questions, got {len(qs)}")
    prefix = "agent-daily-" + d.get("id", "")[6:].replace("-", "") if daily else "agent-" + d.get("id", "")[:2]
    for n, q in enumerate(qs, 1):
        tag = q.get("id", f"#{n}")
        if list(q.keys()) != KEYS:
            errs.append(f"{tag}: keys must be exactly {KEYS} in that order, got {list(q.keys())}")
            continue
        if q["id"] != f"{prefix}-{n}":
            errs.append(f"{tag}: id should be {prefix}-{n}")
        if q["difficulty"] not in ("Medium", "Hard", "Senior"):
            errs.append(f"{tag}: bad difficulty {q['difficulty']!r}")
        for k in KEYS:
            if q[k] in ("", None, []):
                errs.append(f"{tag}: empty {k}")
        for en, zh in (("question", "question_zh"), ("answer", "answer_zh")):
            a, b = shape(blocks(q[en])), shape(blocks(q[zh]))
            if a != b:
                errs.append(f"{tag}: {en} vs {zh} block structure differs:\n    en={a}\n    zh={b}")
            if q[en].count("**") != q[zh].count("**"):
                errs.append(f"{tag}: {en}/{zh} bold markers differ ({q[en].count('**')} vs {q[zh].count('**')})")
            for k in (en, zh):
                for p in unsupported(q[k]):
                    errs.append(f"{tag}: {k}: {p}")
                if "\n\n\n" in q[k]:
                    errs.append(f"{tag}: {k}: triple newline")
        if not (2 <= len(q["followups"]) <= 3) or len(q["followups"]) != len(q["followups_zh"]):
            errs.append(f"{tag}: followups must be 2-3 and match followups_zh")
        if not (3 <= len(q["vocab"]) <= 5) or any(set(v) != {"en", "zh"} or not v["en"] or not v["zh"] for v in q["vocab"]):
            errs.append(f"{tag}: vocab must be 3-5 objects with non-empty en and zh")
        if len(q["answer"]) < 1500:
            errs.append(f"{tag}: answer looks too short ({len(q['answer'])} chars) for a crash-course answer")
    if errs:
        print(f"FAIL {path}")
        for e in errs:
            print("  -", e)
        sys.exit(1)
    print(f"OK {path}: {len(qs)} questions, difficulties {[q['difficulty'] for q in qs]}")


main()
