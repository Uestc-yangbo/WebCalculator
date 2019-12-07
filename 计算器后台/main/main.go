package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/websocket"
)

type node struct {
	value interface{}
	next  *node
}

//Stack 栈的链式结构实现
type Stack struct {
	top    *node
	length int
}

//NewStack 创建一个新栈
func NewStack() *Stack {
	return &Stack{nil, 0}
}

//Len 获取长度
func (s *Stack) Len() int {
	return s.length
}

//Peek 获取栈顶元素
func (s *Stack) Peek() interface{} {
	if s.length == 0 {
		return nil
	}
	return s.top.value
}

//Pop 弹出栈顶元素
func (s *Stack) Pop() interface{} {
	if s.length == 0 {
		return nil
	}
	n := s.top
	s.top = n.next
	s.length--
	return n.value
}

//Push 将元素压入栈
func (s *Stack) Push(value interface{}) {
	n := &node{value, s.top}
	// 从头部插入
	s.top = n
	s.length++
}

//将数字和操作符转为字符串数组
func toExp(str string) []string {
	s := make([]string, 0)
	var t bytes.Buffer
	n := 0 // 用于判断括号是否成对
	for _, r := range str {
		if r == ' ' {
			// 去掉空格
			continue
		}
		if isDigit(r) {
			// 是数字 就写到缓存中
			t.WriteRune(r)
		} else {
			rs := string(r)
			if !isSign(rs) {
				panic("unknown sign: " + rs)
			}
			if t.Len() > 0 {
				// 遇到符号 把缓存中的数字 输出为数
				s = append(s, t.String())
				t.Reset()
			}
			s = append(s, rs)
			if r == '(' {
				n++
			} else if r == ')' {
				n--
			}
		}
	}
	if t.Len() > 0 {
		// 最后一个操作符后面的数字 如果最后一个操作符是 ")" 那么 t.Len() 为0
		s = append(s, t.String())
	}
	if n != 0 {
		panic("the number of '(' is not equal to the number of ')' ")
	}
	return s
}

func printExp(exp []string) {
	for _, s := range exp {
		fmt.Print(s, " ")
	}
	fmt.Println()
}

// 判断是否为数字，这里将小数点也看作数字
func isDigit(r rune) bool {
	if r >= '0' && r <= '9' || r == '.' {
		return true
	}
	return false
}

// 判断是否为符号包括新命名的操作符
func isSign(s string) bool {
	switch s {
	case "+", "-", "*", "/", "(", ")", "g", "n", "s", "c", "t", "p", "k", "d", "!":
		return true
	default:
		return false
	}
}

// 中缀表达式转后缀表达式
func toPostfix(exp []string) []string {
	result := make([]string, 0)
	s := NewStack()
	for _, str := range exp {
		if isSign(str) {
			// 若是符号
			if str == "(" || s.Len() == 0 {
				// "(" 或者 栈为空 直接进栈
				// 括号中的计算 需要单独处理 相当于一个新的上下文
				// 如果栈为空 需要先进栈 和后续操作符比较优先级之后 才能决定计算顺序
				s.Push(str)
			} else {
				if str == ")" {
					// 若为 ")" 依次弹出栈顶元素并输出 直到遇到 "("
					for s.Len() > 0 {
						if s.Peek().(string) == "(" {
							s.Pop()
							break
						}
						result = appendStr(result, s.Pop().(string))
					}
				} else {
					// 判断其与栈顶符号的优先级
					// 如果栈顶是 "(" 说明是新的上下文 不能相互比较优先级
					for s.Len() > 0 && s.Peek().(string) != "(" && signCompare(str, s.Peek().(string)) <= 0 {
						// 当前符号的优先级 不大于栈顶元素 弹出栈顶元素并输出
						// 优先级高的操作 需要先计算
						// 优先级相同 因为栈中的操作是先放进去的 也需要先计算
						result = appendStr(result, s.Pop().(string))
					}
					// 当前符号入栈
					s.Push(str)
				}
			}
		} else {
			// 若是数字就输出
			result = appendStr(result, str)
		}
	}
	for s.Len() > 0 {
		result = appendStr(result, s.Pop().(string))
	}
	return result
}

func appendStr(slice []string, str string) []string {
	if str == "(" || str == ")" {
		// 后缀表达式 不包含括号 这里删除括号
		return slice
	}
	return append(slice, str)
}

// 比较符号优先级
func signCompare(a, b string) int {
	return getSignValue(a) - getSignValue(b)
}

// 优先级越高 值越大
func getSignValue(a string) int {
	switch a {
	case "(", ")":
		return 3
	case "g", "n", "s", "c", "t", "p", "k", "d", "!":
		return 2
	case "*", "/":
		return 1
	default:
		return 0
	}
}

// 通过后缀表达式 计算值
func calValue(exp []string) float64 {
	s := NewStack()
	for _, str := range exp {
		if isSign(str) {
			/*若操作符需要两个操作数，则弹出两个操作数
			计算出结果后入栈，由于栈先进先出的原因
			这里先弹出b，再弹出a*/
			/*如果操作符只需一个操作数，则弹出一个操作数
			计算出结果后入栈*/
			var n float64
			switch str {
			case "+": //相加
				b := getfloat64(s)
				a := getfloat64(s)
				n = a + b
			case "-": //相减
				b := getfloat64(s)
				a := getfloat64(s)
				n = a - b
			case "*": //相乘
				b := getfloat64(s)
				a := getfloat64(s)
				n = a * b
			case "/": //相除
				b := getfloat64(s)
				a := getfloat64(s)
				n = a / b
			case "g": //取对数lg
				b := getfloat64(s)
				n = math.Log10(b)
			case "n": //取自然对数ln
				b := getfloat64(s)
				n = math.Log(b)
			case "s": //三角函数Sin
				b := getfloat64(s)
				//这里的Sin计算的是弧度，需进行转换，下同
				n = math.Sin(b * math.Pi / 180.0)
			case "c": //三角函数Cos
				b := getfloat64(s)
				n = math.Cos(b * math.Pi / 180.0)
			case "t": //三角函数Tan
				b := getfloat64(s)
				n = math.Tan(b * math.Pi / 180.0)
			case "p": //平方操作
				b := getfloat64(s)
				n = b * b
			case "k": //开方操作
				b := getfloat64(s)
				n = math.Sqrt(b)
			case "d": //取倒数
				b := getfloat64(s)
				n = 1 / b
			case "!": //阶乘
				b := getfloat64(s)
				n = 1
				//这里兼容了用户的错误操作，小数向下取整进行阶乘
				for i := 1.0; i <= b; i = i + 1.0 {
					n = n * i
				}
			}
			// 计算结果压栈
			s.Push(n)
		} else {
			// 数字直接压栈
			s.Push(str)
		}
	}
	// 栈顶元素 为最终结果
	return getfloat64(s)
}

// 弹出栈顶元素 并转为float64
func getfloat64(s *Stack) float64 {
	v := s.Pop()
	switch v.(type) {
	case float64: // push进去的计算结果为float64
		return v.(float64)
	case string: // exp中的数据为string
		if i, err := strconv.ParseFloat(v.(string), 64); err != nil {
			panic(err)
		} else {
			return i
		}
	}
	panic(fmt.Sprintf("unknown value type: %T", v))
}

//GetReady 预处理字符串 更换操作符位置
func GetReady(str string) string {
	mytarget := 0
	num1 := 0
	num2 := 0
	for i := 0; i < len(str); i++ {
		str1 := make([]string, 0)
		if str[i] == 'p' || str[i] == 'd' || str[i] == '!' {
			if str[i-1] == ')' { //用于处理含有“（”的情况
				for j := i - 1; j >= 0; j-- {
					if str[j] == ')' {
						num1 = num1 + 1
					}
					if str[j] == '(' {
						num2 = num2 + 1
					}
					if num1 == num2 {
						mytarget = j
						break
					}
				}
			} else { //用于判断没有括号的情况
				for j := i - 1; j >= 0; j-- {
					if !(str[j] >= '0' && str[j] <= '9' || str[j] == '.') {
						mytarget = j + 1
						break
					}
				}
			}
			num1 = 0
			num2 = 0

			if mytarget != 0 {
				str1 = append(str1, str[:mytarget])
			}
			str1 = append(str1, str[i:i+1])
			str1 = append(str1, str[mytarget:i])
			if i != len(str)-1 {
				str1 = append(str1, str[i+1:])
			}
			Str := strings.Join(str1, " ")
			str = strings.Replace(str, str, Str, 1)
		}
	}
	return str
}

//Echo 和前端进行数据交互
func Echo(ws *websocket.Conn) {
	var err error
	for {
		var reply string
		//websocket接受信息
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("receive failed:", err)
			break
		}
		fmt.Println("reveived from client: " + reply)
		//对接收到的字符串进行预处理
		mystr := GetReady(reply)
		fmt.Println(mystr)
		// 转为字符串数组
		exp := toExp(mystr)
		printExp(exp)
		// 转为后缀表达式
		postfixExp := toPostfix(exp)
		fmt.Print("Postfix expression: ")
		printExp(postfixExp)
		// 计算结果
		value := calValue(postfixExp)
		fmt.Println(fmt.Sprintf("Result: %f", value))

		msg := strconv.FormatFloat(value, 'f', -1, 64)
		fmt.Println("send to client:" + msg)
		//这里是发送消息
		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("send failed:", err)
			break
		}
	}
}

func main() {
	//接受websocket的路由地址
	http.Handle("/websocket", websocket.Handler(Echo))
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
