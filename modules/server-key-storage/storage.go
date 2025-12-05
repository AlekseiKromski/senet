package server_key_storage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/AlekseiKromski/server-core/core"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

type ServerKeyStorage interface {
	Encrypt(payload string) (string, error)
	Decrypt(payload string) (string, error)
}

type Config struct {
	publicKeyLocation  string
	privateKeyLocation string
}

func NewConfig(publicKeyLocation, privateKeyLocation string) *Config {
	return &Config{
		publicKeyLocation:  publicKeyLocation,
		privateKeyLocation: privateKeyLocation,
	}
}

type Storage struct {
	config     *Config
	privateKey []byte
	publicKey  []byte

	core.SignedLogger
}

func NewStorage(config *Config) *Storage {
	s := &Storage{
		config: config,
	}

	baseLogger := core.NewDefaultLogger(s.Signature())
	s.SignedLogger = core.NewDefaultSignedLogger(baseLogger)

	return s
}

func (s *Storage) Start(notifyChannel chan struct{}, _ chan core.BusEvent, _ map[string]core.Module) {
	publicKey, err := os.ReadFile(s.config.publicKeyLocation)
	if err != nil {
		s.Log("cannot read public key file", err.Error())
		return
	}

	privateKey, err := os.ReadFile(s.config.privateKeyLocation)
	if err != nil {
		s.Log("cannot read private key file", err.Error())
		return
	}

	s.publicKey = publicKey
	s.privateKey = privateKey

	notifyChannel <- struct{}{}
}

func (p *Storage) Stop() {}

func (p *Storage) Require() []string {
	return []string{}
}

func (p *Storage) Signature() string {
	return "server-key-storage"
}

func (s *Storage) Log(messages ...string) {
	logString := fmt.Sprintf("%s: ", s.Signature())

	for index, message := range messages {
		logString += fmt.Sprintf("%s", message)

		if index != len(messages)-1 {
			logString += ": "
		}
	}
	log.Printf(logString)
}

func (s *Storage) Encrypt(payload string) (string, error) {
	r := bytes.NewReader(s.publicKey)
	entityList, err := openpgp.ReadArmoredKeyRing(r)
	if err != nil {
		s.Log("cannot read armored public key ring", err.Error())
		return "", nil
	}

	var encryptedBuf bytes.Buffer
	w, err := armor.Encode(&encryptedBuf, "PGP MESSAGE", nil)
	if err != nil {
		s.Log("cannot prepare the buffer to write the encrypted message", err.Error())
		return "", nil
	}

	// Encrypt the message
	pt, err := openpgp.Encrypt(w, entityList, nil, nil, nil)
	if err != nil {
		s.Log("cannot encrypt message", err.Error())
		return "", nil
	}
	_, err = pt.Write([]byte(payload))
	if err != nil {
		s.Log("cannot encrypt message", err.Error())
		return "", nil
	}
	pt.Close()
	w.Close()

	return encryptedBuf.String(), nil
}

func (s *Storage) Decrypt(payload string) (string, error) {
	r := bytes.NewReader(s.privateKey)
	entityList, err := openpgp.ReadArmoredKeyRing(r)
	if err != nil {
		s.Log("Error reading armored key ring:", err.Error())
		return "", nil
	}

	// Decrypt the private key using the passphrase
	for _, entity := range entityList {
		if entity.PrivateKey != nil && entity.PrivateKey.Encrypted {
			err := entity.PrivateKey.Decrypt([]byte("senet"))
			if err != nil {
				s.Log("Error decrypting private key:", err.Error())
				return "", nil
			}
		}
		for _, subkey := range entity.Subkeys {
			if subkey.PrivateKey != nil && subkey.PrivateKey.Encrypted {
				err := subkey.PrivateKey.Decrypt([]byte("senet"))
				if err != nil {
					s.Log("Error decrypting subkey:", err.Error())
					return "", nil
				}
			}
		}
	}

	// Decode the armored encrypted message
	armorBlock, err := armor.Decode(bytes.NewReader([]byte(payload)))
	if err != nil {
		s.Log("Error decoding armored message:", err.Error())
		return "", nil
	}

	md, err := openpgp.ReadMessage(armorBlock.Body, entityList, nil, nil)
	if err != nil {
		s.Log("Error reading message:", err.Error())
		return "", nil
	}

	// Decrypt the message
	decryptedBytes, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		s.Log("Error reading decrypted message:", err.Error())
		return "", nil
	}

	return string(decryptedBytes), nil
}
