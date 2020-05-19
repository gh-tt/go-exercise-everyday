package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

func main() {
	host := "smtp.exmail.qq.com"
	port := 465
	email := "zhangyutao@x71x.com"    // 发送邮箱
	pwd := "WgEJ6YBahTu7fAUc"         // 邮箱密码
	toEmail := "519840013@qq.com" // 目标地址
	header := make(map[string]string)
	header["From"] = "涛" + "<" + email + ">"
	header["To"] = toEmail
	header["Subject"] = "abc邮件标题12344"
	header["Content-Type"] = "text/html;charset=UTF-8"
	text := "您好！！！！"
	body := `
    <html>
    <body>
    <h3>
    "测试通过工具发送邮件"` + text + `
    </h3>
    </body>
    </html>
    `
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s:%s\r\n", k, v)
	}
	message += "\r\n" + body
	auth := smtp.PlainAuth(
		"",
		email,
		pwd,
		host,
	)
	fmt.Println("send email")
	err := SendMailUsingTLS(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		email,
		toEmail,
		[]byte(message),
	)
	if err != nil {
		fmt.Println("发送邮件失败!")
		fmt.Println(err)
	} else {
		fmt.Println("发送邮件成功!")
	}
}
func SendMailUsingTLS(addr string, auth smtp.Auth, from string, to string, msg []byte) (err error) {
	c, err := Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	tos := strings.Split(to, ";")
	for _, addr := range tos {
		if err = c.Rcpt(addr); err != nil {
			fmt.Print(err)
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}
