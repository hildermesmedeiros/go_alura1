package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var clear map[string]func() //create a map for storing clear funcs
const monitoramentos = 3    //To monitor 3 times
const delay = 5

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func main() {
	CallClear()
	fmt.Println(time.Now().Format("02/01/2006 15:04:05"))
	exibeIntroducao()
	for {
		exibeMenu()
		comando := leComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			imprimeLogs()
		case 0:
			fmt.Println("Saindo do programa")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando")
			time.Sleep(2 * time.Second)
			CallClear()
		}
	}

}

func exibeIntroducao() {
	nome := "Hildermes"
	versao := 1.1
	fmt.Println("Olá, sr.", nome)
	fmt.Println("Este programa está na versão", versao)
}

func exibeMenu() {
	fmt.Println("1- Iniciar Monitoramento")
	fmt.Println("2- Exibir Logs")
	fmt.Println("0- Sair do Programa")
}

func leComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("O comando escolhido foi", comandoLido)

	return comandoLido
}

func iniciarMonitoramento() {
	CallClear()
	fmt.Println("Monitorando...")
	sites := leSitesDoArquivo()

	for i := 0; i < monitoramentos; i++ {
		for _, site := range sites {
			resp, err := http.Get(site)

			if err != nil {
				fmt.Println("Get Error: ", err)
			}

			if resp.StatusCode == 200 {
				sucesso := textColor("g", "sucesso!")
				fmt.Println("Site:", site, " foi carregado com ", sucesso)
				registraLog(site, true)
			} else {
				failure := textColor("r", strconv.Itoa(resp.StatusCode))
				fmt.Println("Site:", site, "está com problemas. Status Code:", failure)
				registraLog(site, false)
			}
		}
		time.Sleep(delay * time.Second)
		fmt.Println("")
	}
	fmt.Println("")
	time.Sleep(2 * time.Second)
	CallClear()

}

func textColor(colorName string, text string) string {
	colorNameLowerCase := strings.ToLower(colorName)
	color := "\033[0m"
	endColor := "\033[0m"
	switch colorNameLowerCase {
	case "r":
		color = "\033[31m"
		return string(color) + text + endColor
	case "g":
		color = "\033[32m"
		return string(color) + text + endColor
	case "y":
		color = "\033[33m"
		return string(color) + text + endColor
	case "b":
		color = "\033[34m"
		return string(color) + text + endColor
	case "p":
		color = "\033[35m"
		return string(color) + text + endColor
	case "c":
		color = "\033[36m"
		return string(color) + text + endColor
	case "w":
		color = "\033[36m"
		return string(color) + text + endColor
	default:
		color = "\033[0m"
		return string(color) + text
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! The script can't clear terminal screen!!")
	}
}

// restante do código omitido

func leSitesDoArquivo() []string {

	var sites []string

	arquivo, err := os.Open("sites.txt")
	if err != nil {
		fmt.Println("Open Error: ", err)
		os.Exit(-1)
	}

	leitor := bufio.NewReader(arquivo)

	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)
		sites = append(sites, linha)
		if err == io.EOF {
			break
		}
	}

	arquivo.Close()

	return sites
}

func registraLog(site string, status bool) {

	arquivo, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("OpenFile Error: ", err)
	}

	if status == true {
		arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")
	} else {
		arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")
	}

	arquivo.Close()
}

func imprimeLogs() {
	CallClear()
	arquivo, err := ioutil.ReadFile("log.txt")

	if err != nil {
		fmt.Println("ReadFile Error: 1", err)
	}
	arquivo = []byte(strings.ReplaceAll(string(arquivo), "true", textColor("g", "true")))
	arquivo = []byte(strings.ReplaceAll(string(arquivo), "false", textColor("r", "false")))
	fmt.Println(string(arquivo))
}
