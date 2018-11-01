package mypackages

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func SSHExec(host string, command string) (string, error) {
	key, err := ioutil.ReadFile("./src/id_rsa") // For local sistem
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
		return "error1", err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
		return "error2", err
	}
	// Конфигурация для подключения по ssh
	sshConfig := &ssh.ClientConfig{
		User: os.Getenv("LOGNAME"),
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		//TODO: Для продакшена не использовать режим Insecure, реализовать механизм определения окружений и подстановки значений
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	// Инициализируем клиента для подключения к хосту
	client, err := ssh.Dial("tcp", host+":22", sshConfig)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
		return "error3", err
	}

	// Создаем новую сессию
	session, err := client.NewSession()
	if err != nil {
		// log.Fatal("Failed to create session: ", err)
		return "error4", err
	}
	defer session.Close()

	// Выполняем команду на удаленном хосте
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		// log.Fatal("Failed to run: " + err.Error())
		return "error5", err
	}

	return b.String(), err
}
