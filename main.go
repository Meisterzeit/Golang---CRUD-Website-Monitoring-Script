package main

import (
	"bufio" //biblioteca que permite a leitura, inputs e outputs de arquivos externos
	"fmt"
	"io"
	"io/ioutil"
	"net/http" //biblioteca que permite a manipulação de protocolos http
	"os"       //biblioteca que permite comunicacao com o sistema operacional (informa tipos de return)
	"strconv"
	"strings"
	"time"
)

//constante global para o numero de testes feitos na função de monitoramento
const monitoramentos = 2

//constante global para o numero de segundos da função sleep no interior da função monitoramento
const delay = 3

//função sem retorno que solicita do usuário os dados pessoais e dá as boas vindas
func exibeIntroducao() {
	var nome string
	versao := 1.1

	fmt.Println("Insira seu nome: ")
	fmt.Scan(&nome)

	fmt.Println("Bem vindo sr(a) ", nome)
	fmt.Println("A versão deste programa é:", versao)
}

//função sem retorno que exibe o menu de opções para o usuário
func exibeMenu() {

	fmt.Println("Menu de comandos:")
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir Logs")
	fmt.Println("0 - Sair do programa")

}

//função com retorno que recebe o input do usuário para iniciar uma função de comando
func leComando() int {

	var comandoLido int

	fmt.Scan(&comandoLido)
	fmt.Println("O comando escolhido foi:", comandoLido)

	return comandoLido
}

//função que recebe da func leComando() o comando escolhido e retorna funcões das operações propostas
func executaOperacao() int {
	comando := leComando()

	//switch de opções que recebe da váriavel comando o input do usuário
	switch comando {

	case 0:
		fmt.Println("Saindo do programa")
		//func da lib os que sai do programa com sucesso com o valor 0(zero)
		os.Exit(0)
	case 1:
		iniciarMonitoramento()
	case 2:
		imprimeLogs()
	default:
		fmt.Println("Comando inválido, por favor insira comando listado no menu")
		//func da lib os que sai do programa apresentando erro (-1: exit status 255)
		os.Exit(-1)
	}

	return comando
}

func iniciarMonitoramento() {
	fmt.Println("Iniciando monitorando...")
	//slice de sites para serem testados seus status
	//sites := []string{"https://random-status-code.herokuapp.com", "https://www.github.com", "https://www.google.com"}

	sites := leSitesDoArquivo()

	// loop para executar o teste para os sites de acordo com a constante global de monitoramento (linha 11)
	for i := 0; i < monitoramentos; i++ {
		//estrutura de for em Go onde index representa o que irá iterar, site é a variável que será preenchida em cada iteração
		//range simboliza o tamanho do slice(chamado de sites) a ser percorrido.
		//Igual a for index = 0; index < len(sites); index++{...}
		//Enquanto existirem elementos no slice, o range irá iterar até o último e acrescentar +1 a variável index e ler o elemento contido na posicão
		for index, site := range sites {
			fmt.Println("Testando site:", index, ":", site)
			testaSite(site)
		}
		//função da biblioteca time para intervalar cada teste a cada vez que a const delay permitir (de acordo com o loop externo).
		time.Sleep(delay * time.Second)

		//pular uma linha a cada execução do programa no terminal
		fmt.Println("")
	}
	//pular uma linha a cada execução do programa no terminal
	fmt.Println("")
}

func testaSite(site string) {

	//função da biblioteca http que retorna um header a partir de solicitação GET (http.Get)
	//um under line "_" pode ser inserido para ignorar um dos parametros da funcão, uma vez que o http.Get solicita um segundo
	//parametro de tratamento de erros.
	resp, err := http.Get(site)

	//trata erros que possam ocorrer na leitura (http.Get) da variável "sites"
	if err != nil {
		fmt.Println("Ocorreu um erro na leitura", err)
	}

	//Status como é uma função da biblioteca http que permite acessar o código de resposta de uma requisição
	if resp.StatusCode == 200 {
		fmt.Println("O site ", site, "foi carregado com sucesso e o StatusCode:", resp.StatusCode)
		//chama a função registraLog() para coletar informações do request do site caso o status seja 200 (true)
		registraLog(site, true)
	} else {
		//qualquer statusCode != 200, o loop exibe o código de erro
		fmt.Println("O site ", site, "não foi carregado com sucesso e o StatusCode:", resp.StatusCode)
		//chama a função registraLog() para coletar informações do request do site caso o status seja != 200 (false)
		registraLog(site, false)
	}

}

//função com retorno que irá receber um slice de strings proveniente de um txt
func leSitesDoArquivo() []string {

	var sites []string

	//variavel que recebe o txt externo a partir da func os.Open
	//esta recebe dois parametros, sendo o segundo um tratamento de erro
	arquivo, err := os.Open("sites.txt")
	//trata erros que possam ocorrer na leitura (os.Open) do arquivo "sites.txt"
	if err != nil {
		fmt.Println("Ocorreu um erro na leitura", err)
	}

	//bufio.NewReader(arquivo) e uma função que fará a leitura, linha a linha do arquivo externo
	leitor := bufio.NewReader(arquivo)
	//loop infinito que irá executar o script fazendo a leitura de cada linha do txt
	for {

		//ReaderString(arquivo) é uma função que fará a leitura do arquivo, limitando ao último caracter inserido na forma de bytes
		//Ela é armazenada em uma variável e trata possíveis erros
		linha, err := leitor.ReadString('\n')
		//TrimSpace é uma função da biblioteca strings que irá cortar excessos de espaços na impressão do arquivo lido
		linha = strings.TrimSpace(linha)
		//o append irá inserir no slice  "sites" cada linha do txt no log do arquivo lido
		sites = append(sites, linha)
		//tratamento de erro caso ocorra erro ao chegar ao fim do arquivo (End Of File)
		if err == io.EOF {
			//break para sair do loop infinito após alcançar o fim do txt
			break
		}

	}
	//Função da biblioteca os para encerrar o arquivo
	arquivo.Close()

	return sites
}

func registraLog(site string, status bool) {
	//os.OpenFile para abrir um arquivo txt com as infos de log da requisição
	//ela recebe 3 parametro: o arquivo txt, a flag(vide documentação golang) com a opção de ler e escrever no arquivo txt
	// ou caso o txt não exista, a opção de criá-lo, por fim os.O_APPEND irá escrever cada linha de log.
	// Por fim recebe as permissões (função 0666).
	arquivo, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	//tratamento de error
	if err != nil {
		fmt.Println("Ocorreu um erro de log", err)
	}
	//função writeString da biblioteca String para escrever no log o que deve ser impresso
	//time.Now().Format() é uma função da biblioteca time para inserir tempo com a data atual (vide doc do Golang para a formatação)
	//ela vai receber uma concatenação que terá a func da biblioteca strConv.FormatBool para converter o boolean para string
	//"\n" ao fim para separar os logs por linhas no arquivo txt.
	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " _ " + site + "- online: " + strconv.FormatBool(status) + "\n")

	arquivo.Close()

}

func imprimeLogs() {
	fmt.Println("Exibindo logs do programa")
	//função para ler o arquivo de log (irá abrir, ler e fechar o arquivo != da lib "os" que trabalha em baixo nivel)
	arquivo, err := ioutil.ReadFile("log.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro de log:", err)
	}
	//irá imprimir o log no formato string, sem esse tratamento o output seria em formato de bytes

	fmt.Println(string(arquivo))

}

func main() {

	exibeIntroducao()

	//Irá criar um loop infinito para executar o script até a opção 0 ser colocada como input
	for {
		exibeMenu()
		executaOperacao()
	}

}
