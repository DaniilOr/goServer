package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)
var (
	errWrongRequest = errors.New("Worng request")
	errFormatMissmatch = errors.New("Wrong file format")
)
func main() {
	err := execute()
	if err != nil{
		os.Exit(1)
	}
}
func execute() (err error){
	listener, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil{
		log.Println(err)
		return err
	}
	defer func(){
		if cerr := listener.Close(); cerr != nil{
			log.Println(cerr)
			if err == nil{
				err = cerr
			}
		}
	}()
	for{
		conn, err := listener.Accept()
		if err != nil{
			log.Println(err)
			return  err
		}
		handle(conn)
	}
}

func handle(conn net.Conn){
	defer func(){
		if cerr := conn.Close(); cerr != nil{
			log.Println(cerr)
		}
	}()
	reader := bufio.NewReader(conn)
	const delim = '\n'
	line, err := reader.ReadString(delim)
	/*for{

		if err != nil{
			if err != io.EOF{
				log.Println(err)
			}
			log.Printf("recieved: %s\n", line)
			return
		}
		log.Printf("recieved: %s\n", line)
	}*/
	parts := strings.Split(line, " ")
	if len(parts) != 3{
		log.Println()
	}

	if err != nil{
		log.Println(errWrongRequest)
		return
	}

	path := parts[1]
	switch {
	case path == "/":
		err = WriteIndex(conn)
		break
	case strings.Contains(path,"/operations"):
		err = WriteOperations(conn, path)
		break
	default:
		err = Write404(conn)
	}
	if err != nil{
		log.Println(err)
	}
}

func WriteIndex(writer io.Writer) (err error){
	page, err := ioutil.ReadFile("web/index.html")
	if err != nil{
		return err
	}
	username := "User"
	balance := "1_000_000"
	page = bytes.ReplaceAll(page, []byte("{user}"), []byte(username))
	page = bytes.ReplaceAll(page, []byte("{balance}"), []byte(balance))
	return WriteResponse(writer, 200, []string{
		"Content-type: text/html;charset=utf-8",
		fmt.Sprintf("Content-length: %d", len(page)),
		"Conntection: close",}, page)

}
func WriteResponse(writer io.Writer, status int, headers []string, content []byte) (err error){
	const CRLF = "\r\n"
	w := bufio.NewWriter(writer)
	_, err = w.WriteString(fmt.Sprintf("HTTP/1.1 %d OK", status) + CRLF)
	if err != nil{
		return err
	}
	for _, h := range headers{
		_, err = w.WriteString(h + CRLF)
		if err != nil{
			return err
		}
	}
	_, err = w.WriteString(CRLF)
	if err != nil{
		return err
	}
	_, err = w.Write(content)
	if err != nil{
		return err
	}
	err = w.Flush()
	if err != nil{
		return err
	}
	return nil
}
func Write404(writer io.Writer) (err error){
	page, err := ioutil.ReadFile("web/404.html")
	if err != nil{
		return nil
	}
	return WriteResponse(writer, 200, []string{
		"Content-type: text/html;charset=utf-8",
		fmt.Sprintf("Content-length: %d", len(page)),
		"Conntection: close",}, page)


}
func WriteOperations(writer io.Writer, command string) (err error){
	format := strings.Split(command, ".")[1]
	switch format{
	case "csv":
		file, err := ioutil.ReadFile("operations.csv")
		if err != nil{
			return err
		}
		return WriteResponse(writer, 200, []string{
			"Content-type: text; charset=utf-8",
			fmt.Sprintf("Content-length: %d", len(file)),
			"Conntection: close",}, file)
	case "json":
		file, err := ioutil.ReadFile("operations.json")
		if err != nil{
			return err
		}
		return WriteResponse(writer, 200, []string{
			"Content-type: application/json; charset=utf-8",
			fmt.Sprintf("Content-length: %d", len(file)),
			"Conntection: close",}, file)

	case "xml":
		file, err := ioutil.ReadFile("operations.xml")
		if err != nil{
			return err
		}
		return WriteResponse(writer, 200, []string{
			"Content-type: application/xml; charset=utf-8",
			fmt.Sprintf("Content-length: %d", len(file)),
			"Conntection: close",}, file)

	default:
		return errFormatMissmatch
	}

}