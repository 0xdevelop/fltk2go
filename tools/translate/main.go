package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// 获取所有头文件
	includeDir := "libs/fltk/include_cn"
	headerFiles, err := getHeaderFiles(includeDir)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("找到 %d 个头文件\n", len(headerFiles))

	// 翻译每个文件
	for _, filePath := range headerFiles {
		translateFile(filePath)
	}

	fmt.Println("所有文件翻译完成！")
}

// getHeaderFiles 获取所有头文件
func getHeaderFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := filepath.Ext(path)
			if strings.ToLower(ext) == ".h" || ext == ".H" {
				files = append(files, path)
			}
		}

		return nil
	})

	return files, err
}

// translateFile 翻译单个文件的所有注释
func translateFile(filePath string) {
	fmt.Printf("正在翻译: %s\n", filePath)

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("读取文件错误: %v\n", err)
		return
	}

	// 转换为字符串
	text := string(content)

	// 处理多行注释 /** */
	multilineCommentRegex := regexp.MustCompile(`/\*\*[\s\S]*?\*/`)
	text = multilineCommentRegex.ReplaceAllStringFunc(text, translateComment)

	// 处理单行注释 //
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		// 检查是否在字符串中
		inString := false
		stringChar := ""
		for j, c := range line {
			if c == '"' || c == '\'' {
				// 检查是否是转义的引号
				if j == 0 || line[j-1] != '\\' {
					if !inString {
						inString = true
						stringChar = string(c)
					} else if string(c) == stringChar {
						inString = false
					}
				}
			}
			// 检查是否是单行注释
			if c == '/' && j+1 < len(line) && line[j+1] == '/' && !inString {
				// 找到单行注释
				comment := line[j:]
				translatedComment := translateComment(comment)
				lines[i] = line[:j] + translatedComment
				break
			}
		}
	}

	// 重新组合内容
	translatedContent := strings.Join(lines, "\n")

	// 写回文件
	err = os.WriteFile(filePath, []byte(translatedContent), 0644)
	if err != nil {
		fmt.Printf("写入文件错误: %v\n", err)
		return
	}

	fmt.Printf("翻译完成: %s\n", filePath)
}

// translateComment 翻译单个注释块
func translateComment(comment string) string {
	// 基础翻译
	comment = regexp.MustCompile(`// Main header file for the Fast Light Tool Kit \(FLTK\).`).ReplaceAllString(comment, "// Fast Light Tool Kit (FLTK) 的主头文件。")
	comment = regexp.MustCompile(`// Copyright 1998-2024 by Bill Spitzak and others.`).ReplaceAllString(comment, "// 版权所有 1998-2024 Bill Spitzak 及其他贡献者。")
	comment = regexp.MustCompile(`// This library is free software. Distribution and use rights are outlined in`).ReplaceAllString(comment, "// 本库是自由软件。分发和使用权利在以下文件中概述：")
	comment = regexp.MustCompile(`// the file "COPYING" which should have been included with this file.  If this`).ReplaceAllString(comment, "// 文件 \"COPYING\" 应随本文件一起包含。如果此")
	comment = regexp.MustCompile(`// file is missing or damaged, see the license at:`).ReplaceAllString(comment, "// 文件丢失或损坏，请查看以下许可证：")
	comment = regexp.MustCompile(`//     https://www.fltk.org/COPYING.php`).ReplaceAllString(comment, "//     https://www.fltk.org/COPYING.php")
	comment = regexp.MustCompile(`// Please see the following page on how to report bugs and issues:`).ReplaceAllString(comment, "// 请查看以下页面了解如何报告错误和问题：")
	comment = regexp.MustCompile(`//     https://www.fltk.org/bugs.php`).ReplaceAllString(comment, "//     https://www.fltk.org/bugs.php")

	// 常用短语翻译
	phrases := map[string]string{
		"Returns the":          "返回",
		"Returns a":            "返回一个",
		"Returns non-zero":     "返回非零值",
		"Returns zero":         "返回零",
		"Returns true":         "返回真",
		"Returns false":        "返回假",
		"This is":              "这是",
		"This function":        "此函数",
		"This method":          "此方法",
		"Use this":             "使用此",
		"You can":              "您可以",
		"You should":           "您应该",
		"You must":             "您必须",
		"If you":               "如果您",
		"is used to":           "用于",
		"can be used to":       "可用于",
		"for use with":         "用于与",
		"in order to":          "为了",
		"such as":              "例如",
		"like":                 "例如",
		"should be":            "应该是",
		"must be":              "必须是",
		"will be":              "将是",
		"can be":               "可以是",
		"cannot be":            "不能是",
		"is not":               "不是",
		"does not":             "不",
		"do not":               "不要",
		"has been":             "已经",
		"have been":            "已经",
		"was":                  "是",
		"were":                 "是",
		"are":                  "是",
		"is":                   "是",
		"the":                  "的",
		"a":                    "一个",
		"an":                   "一个",
	}

	// 替换常用短语
	for en, zh := range phrases {
		comment = strings.ReplaceAll(comment, en, zh)
	}

	// 特定术语翻译
	terms := map[string]string{
		"widget":           "控件",
		"window":           "窗口",
		"event":            "事件",
		"callback":         "回调",
		"handle":           "处理",
		"draw":             "绘制",
		"size":             "大小",
		"width":            "宽度",
		"height":           "高度",
		"position":         "位置",
		"color":            "颜色",
		"font":             "字体",
		"text":             "文本",
		"image":            "图像",
		"pixmap":           "像素图",
		"bitmap":           "位图",
		"surface":          "表面",
		"buffer":           "缓冲区",
		"clipboard":        "剪贴板",
		"selection":        "选择",
		"drag":             "拖动",
		"drop":             "放置",
		"shortcut":         "快捷键",
		"modifier":         "修饰键",
		"keypress":         "按键",
		"mouse":            "鼠标",
		"click":            "点击",
		"double-click":     "双击",
		"button":           "按钮",
		"key":              "键",
		"screen":           "屏幕",
		"display":          "显示器",
		"scale":            "缩放",
		"resolution":       "分辨率",
		"colormap":         "颜色映射表",
		"palette":          "调色板",
		"border":           "边框",
		"shadow":           "阴影",
		"radius":           "半径",
		"rounded":          "圆角",
		"box":              "盒子",
		"boxtype":          "盒子类型",
		"label":            "标签",
		"labeltype":        "标签类型",
		"menu":             "菜单",
		"item":             "项",
		"list":             "列表",
		"browser":          "浏览器",
		"scroll":           "滚动",
		"scrollbar":        "滚动条",
		"tab":              "标签页",
		"panel":            "面板",
		"group":            "组",
		"pack":             "包装",
		"tile":             "瓦片",
		"frame":            "框架",
		"canvas":           "画布",
		"chart":            "图表",
		"progress":         "进度",
		"counter":          "计数器",
		"spinner":          "微调器",
		"slider":           "滑块",
		"dial":             "拨号盘",
		"roller":           "滚轮",
		"positioner":       "定位器",
		"valuator":         "估值器",
		"adjuster":         "调节器",
		"clock":            "时钟",
		"calendar":         "日历",
		"chooser":          "选择器",
		"color chooser":     "颜色选择器",
		"file chooser":      "文件选择器",
		"native":           "原生",
		"export":           "导出",
		"import":           "导入",
		"print":            "打印",
		"printer":          "打印机",
		"handler":          "处理器",
		"driver":           "驱动程序",
		"plugin":           "插件",
		"scheme":           "方案",
		"theme":            "主题",
		"style":            "样式",
		"title":            "标题",
		"icon":             "图标",
		"tooltip":          "工具提示",
		"status":           "状态",
		"error":            "错误",
		"warning":          "警告",
		"debug":            "调试",
		"info":             "信息",
		"log":              "日志",
		"fatal":            "致命",
		"assert":           "断言",
		"abort":            "中止",
		"exit":             "退出",
		"quit":             "退出",
		"return":           "返回",
		"break":            "中断",
		"continue":         "继续",
		"default":          "默认",
		"static":           "静态",
		"const":            "常量",
		"volatile":         "易变",
		"extern":           "外部",
		"inline":           "内联",
		"virtual":          "虚拟",
		"override":         "覆盖",
		"final":            "最终",
		"friend":           "友元",
		"namespace":        "命名空间",
		"template":         "模板",
		"typename":         "类型名",
		"class":            "类",
		"struct":           "结构体",
		"union":            "联合体",
		"enum":             "枚举",
		"typedef":          "类型定义",
		"constructor":      "构造函数",
		"destructor":       "析构函数",
		"copy constructor": "拷贝构造函数",
		"move constructor": "移动构造函数",
		"copy assignment":  "拷贝赋值",
		"move assignment":  "移动赋值",
		"initialization":   "初始化",
		"destruction":      "销毁",
		"allocation":       "分配",
		"deallocation":     "释放",
		"memory":           "内存",
		"heap":             "堆",
		"stack":            "栈",
		"global":           "全局",
		"local":            "局部",
		"member":           "成员",
		"field":            "字段",
		"property":         "属性",
		"method":           "方法",
		"function":         "函数",
		"procedure":        "过程",
		"subroutine":       "子程序",
		"parameter":        "参数",
		"argument":         "实参",
		"return value":     "返回值",
		"error code":       "错误码",
		"success":          "成功",
		"failure":          "失败",
		"result":           "结果",
		"output":           "输出",
		"input":            "输入",
		"input/output":     "输入/输出",
		"interface":        "接口",
		"implementation":   "实现",
		"abstraction":      "抽象",
		"encapsulation":    "封装",
		"inheritance":      "继承",
		"polymorphism":     "多态",
		"overloading":      "重载",
		"overriding":       "覆盖",
		"generic":          "泛型",
		"macro":            "宏",
		"preprocessor":     "预处理器",
		"compiler":         "编译器",
		"linker":           "链接器",
		"assembler":        "汇编器",
		"interpreter":      "解释器",
		"just-in-time":     "即时",
		"optimization":     "优化",
		"debugging":        "调试",
		"profiling":        "性能分析",
		"benchmark":        "基准测试",
		"performance":      "性能",
		"efficiency":       "效率",
		"speed":            "速度",
		"memory usage":     "内存使用",
		"cpu usage":        "CPU使用",
		"resource":         "资源",
		"system":           "系统",
		"platform":         "平台",
		"operating system": "操作系统",
		"console":          "控制台",
		"terminal":         "终端",
		"command":          "命令",
		"option":           "选项",
		"flag":             "标志",
		"switch":           "开关",
		"configuration":    "配置",
		"setting":          "设置",
		"preference":       "首选项",
		"variable":         "变量",
		"constant":         "常量",
		"value":            "值",
		"pointer":          "指针",
		"reference":        "引用",
		"array":            "数组",
		"vector":           "向量",
		"map":              "映射",
		"set":              "集合",
		"queue":            "队列",
		"linked list":      "链表",
		"tree":             "树",
		"graph":            "图",
		"node":             "节点",
		"edge":             "边",
		"leaf":             "叶子",
		"root":             "根",
		"parent":           "父",
		"child":            "子",
		"sibling":          "兄弟",
		"iterator":         "迭代器",
		"algorithm":        "算法",
		"sort":             "排序",
		"search":           "搜索",
		"find":             "查找",
		"insert":           "插入",
		"delete":           "删除",
		"remove":           "移除",
		"update":           "更新",
		"add":              "添加",
		"append":           "追加",
		"prepend":          "前置",
		"concat":           "连接",
		"merge":            "合并",
		"split":            "分割",
		"join":             "连接",
		"substring":        "子字符串",
		"string":           "字符串",
		"char":             "字符",
		"byte":             "字节",
		"bit":              "位",
		"word":             "字",
		"double word":      "双字",
		"quad word":        "四字",
		"octet":            "八位组",
		"signed":           "有符号",
		"unsigned":         "无符号",
		"integer":          "整数",
		"float":            "float",
		"double":           "double",
		"long double":      "long double",
		"bool":             "bool",
		"void":             "void",
		"null":             "NULL",
		"nullptr":          "nullptr",
		"zero":             "零",
		"one":              "一",
		"true":             "真",
		"false":            "假",
	}

	// 替换特定术语
	for en, zh := range terms {
		comment = strings.ReplaceAll(comment, en, zh)
	}

	return comment
}
