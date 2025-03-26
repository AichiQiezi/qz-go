package id

import (
	"errors"
)

// RequestIDGenerator 定义请求 ID 生成器接口
type RequestIDGenerator interface {
	Generate() string
}

// 生成器工厂
var generatorFactory = map[string]func() RequestIDGenerator{
	"default": func() RequestIDGenerator { return &DefaultRequestIDGenerator{} },
	//"custom":  func() RequestIDGenerator { return &CustomRequestIDGenerator{} },
}

// CreateGenerator 创建 ID 生成器
func CreateGenerator(name string) (RequestIDGenerator, error) {
	if factory, exists := generatorFactory[name]; exists {
		return factory(), nil
	}
	return nil, errors.New("unknown generator type")
}

// Holder 单例
type Holder struct {
	requestIDGenerator RequestIDGenerator
}

var instance *Holder

// Init 初始化
func Init(generatorType string) error {
	generator, err := CreateGenerator(generatorType)
	if err != nil {
		return err
	}
	instance = &Holder{requestIDGenerator: generator}
	return nil
}

// GetInstance 获取实例
func GetInstance() *Holder {
	if instance == nil {
		_ = Init("default") // 默认初始化
	}
	return instance
}

// Generate 生成 ID
func (h *Holder) Generate() (string, error) {
	if h.requestIDGenerator == nil {
		return "", errors.New("request ID generator is not initialized")
	}
	return h.requestIDGenerator.Generate(), nil
}

// getRequestIDGeneratorClass 获取 ID 生成器类名（模拟从配置读取）
func getRequestIDGeneratorClass() string {
	return ""
}

// DefaultRequestIDGenerator 默认 ID 生成器
type DefaultRequestIDGenerator struct{}

// Generate 生成默认 ID
func (d *DefaultRequestIDGenerator) Generate() string {
	return "default-request-id"
}
