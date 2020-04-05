package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var commands = map[string]func([]string){
	"help":    printHelp,
	"newkey":  newKey,
	"pubinfo": pubInfo,
	"priinfo": priInfo,
	"gogen":   goGen,
	"getpub":  getPub,
}

var curves = map[string]elliptic.Curve{
	"p224": elliptic.P224(),
	"p256": elliptic.P256(),
	"p384": elliptic.P384(),
	"p521": elliptic.P521(),
}

func main() {
	if len(os.Args) == 1 {
		printHelp(nil)
		return
	}

	command, ok := commands[strings.ToLower(os.Args[1])]
	if !ok {
		fmt.Println("Invalid command", os.Args[1])
		printHelp(nil)
	} else {
		command(os.Args[2:])
	}
}

func printHelp(args []string) {
	fmt.Println("Will's Go ECDSA Key Tool")
	fmt.Println("Commands:")
	fmt.Println("  newkey <curve> <private.der> [public.key]")
	fmt.Println("  getpub <private.der> <public.key>")
	fmt.Println("  pubinfo <public.key>")
	fmt.Println("  priinfo <private.der>")
	fmt.Println("  gogen <public.key> <package> <name> <file>")
}

func newKey(args []string) {
	if len(args) < 2 {
		fmt.Printf("Usage: %s newkey <curve> <private.der> [public.key]", os.Args[0])
	} else {
		curveName := strings.ToLower(args[0])
		privateFile := args[1]
		publicFile := ""
		if len(args) > 2 {
			publicFile = args[2]
		}

		curve, ok := curves[curveName]

		if !ok {
			fmt.Println("Invalid curve", curveName)
			fmt.Println("Curves:")
			for key := range curves {
				fmt.Println(key)
			}
		}

		key, err := ecdsa.GenerateKey(curve, rand.Reader)

		if err != nil {
			fmt.Println("Error generating key")
			fmt.Println(err)
		}

		keyDer, err := x509.MarshalECPrivateKey(key)

		if err != nil {
			fmt.Println("Error marshaling private key")
			fmt.Println(err)
		} else {
			err = ioutil.WriteFile(privateFile, keyDer, 0o500)
			if err != nil {
				fmt.Println("Error writing private key")
				fmt.Println(err)
			}
		}

		if publicFile != "" {
			pubDer, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
			if err != nil {
				fmt.Println("Error marshaling public key")
			} else {
				err = ioutil.WriteFile(publicFile, pubDer, 0o555)
				if err != nil {
					fmt.Println("Error writing public key")
					fmt.Println(err)
				}
			}
		}
	}
}

func pubInfo(args []string) {
	if len(args) < 1 {
		fmt.Printf("Usage: %s pubinfo <public.key>", os.Args[0])
	} else {
		pubDer, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Printf("Error reading public key")
			return
		}

		var pubKey *ecdsa.PublicKey

		pubKeyTemp, err := x509.ParsePKIXPublicKey(pubDer)
		if err != nil {
			fmt.Printf("Error parsing public key")
			return
		}
		pubKey = pubKeyTemp.(*ecdsa.PublicKey)

		fmt.Println("Curve:", pubKey.Curve.Params().Name)
		fmt.Println("X:", pubKey.X)
		fmt.Println("Y:", pubKey.Y)
	}
}

func priInfo(args []string) {
	if len(args) < 1 {
		fmt.Printf("Usage: %s priinfo <private.key>", os.Args[0])
	} else {
		priDer, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Printf("Error reading public key")
			return
		}

		priKey, err := x509.ParseECPrivateKey(priDer)
		if err != nil {
			fmt.Printf("Error parsing public key")
			return
		}

		fmt.Println("Curve:", priKey.Curve.Params().Name)
		fmt.Println("X:", priKey.X)
		fmt.Println("Y:", priKey.Y)
		fmt.Println("D:", priKey.D)
	}
}

func getPub(args []string) {
	if len(args) < 2 {
		fmt.Printf("Usage: %s getpub <private.der> <public.key>", os.Args[0])
	} else {
		priDer, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Printf("Error reading private key")
			return
		}

		priKey, err := x509.ParseECPrivateKey(priDer)
		if err != nil {
			fmt.Printf("Error parsing private key")
			return
		}

		pubDer, err := x509.MarshalPKIXPublicKey(&priKey.PublicKey)
		if err != nil {
			fmt.Println("Error marshaling public key")
			return
		}

		err = ioutil.WriteFile(args[1], pubDer, 0o555)
		if err != nil {
			fmt.Println("Error writing public key")
			return
		}
	}
}

func goGen(args []string) {
	if len(args) < 4 {
		fmt.Printf("Usage: %s gogen <public.key> <package> <name> <file>", os.Args[0])
	} else {
		pubDer, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Printf("Error reading public key")
			return
		}

		var pubKey *ecdsa.PublicKey

		pubKeyTemp, err := x509.ParsePKIXPublicKey(pubDer)
		if err != nil {
			fmt.Printf("Error parsing public key")
			return
		}
		pubKey = pubKeyTemp.(*ecdsa.PublicKey)

		pkg := args[1]
		name := args[2]

		file, err := os.Create(args[3])
		if err != nil {
			fmt.Printf("Error opening file")
			return
		}
		defer file.Close()

		fmt.Fprintf(file, "package %s\n\n", pkg)
		fmt.Fprintf(file, "import (\n\tecdsa \"crypto/ecdsa\"\n\tbig \"math/big\"\n\telliptic \"crypto/elliptic\"\n)\n\n")
		fmt.Fprintf(file, "func %s() *ecdsa.PublicKey {\n", name)
		fmt.Fprintf(file, "\tkey := new(ecdsa.PublicKey)\n")
		fmt.Fprintf(file, "\tkey.Curve = elliptic.%s()\n", strings.ReplaceAll(pubKey.Curve.Params().Name, "-", ""))
		fmt.Fprintf(file, "\tkey.X = big.NewInt(0)\n")
		fmt.Fprintf(file, "\tkey.X.SetBytes(%#v)\n", pubKey.X.Bytes())
		fmt.Fprintf(file, "\tkey.Y = big.NewInt(0)\n")
		fmt.Fprintf(file, "\tkey.Y.SetBytes(%#v)\n", pubKey.Y.Bytes())
		fmt.Fprintf(file, "\treturn key\n")
		fmt.Fprintf(file, "}\n")
	}
}
