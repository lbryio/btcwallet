// Copyright (c) 2015-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

//go:build !js
// +build !js

package prompt

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lbryio/lbcutil/hdkeychain"
	"github.com/lbryio/lbcwallet/internal/legacy/keystore"
	"golang.org/x/crypto/ssh/terminal"
)

// ProvideSeed is used to prompt for the wallet seed which maybe required during
// upgrades.
func ProvideSeed() func() ([]byte, error) {
	return func() ([]byte, error) {
		return provideSeed(bufio.NewReader(os.Stdin))
	}
}
func provideSeed(reader *bufio.Reader) ([]byte, error) {
	for {
		fmt.Print("Enter existing wallet seed: ")
		seedStr, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		seedStr = strings.TrimSpace(strings.ToLower(seedStr))

		seed, err := hex.DecodeString(seedStr)
		if err != nil || len(seed) < hdkeychain.MinSeedBytes ||
			len(seed) > hdkeychain.MaxSeedBytes {

			fmt.Printf("Invalid seed specified.  Must be a "+
				"hexadecimal value that is at least %d bits and "+
				"at most %d bits\n", hdkeychain.MinSeedBytes*8,
				hdkeychain.MaxSeedBytes*8)
			continue
		}

		return seed, nil
	}
}

// ProvidePrivPassphrase is used to prompt for the private passphrase which
// maybe required during upgrades.
func ProvidePrivPassphrase() ([]byte, error) {
	prompt := "Enter the private passphrase of your wallet: "
	for {
		fmt.Print(prompt)
		pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return nil, err
		}
		fmt.Print("\n")
		pass = bytes.TrimSpace(pass)
		if len(pass) == 0 {
			continue
		}

		return pass, nil
	}
}

// promptList prompts the user with the given prefix, list of valid responses,
// and default list entry to use.  The function will repeat the prompt to the
// user until they enter a valid response.
func promptList(reader *bufio.Reader, prefix string, validResponses []string, defaultEntry string) (string, error) {
	// Setup the prompt according to the parameters.
	validStrings := strings.Join(validResponses, "/")
	var prompt string
	if defaultEntry != "" {
		prompt = fmt.Sprintf("%s (%s) [%s]: ", prefix, validStrings,
			defaultEntry)
	} else {
		prompt = fmt.Sprintf("%s (%s): ", prefix, validStrings)
	}

	// Prompt the user until one of the valid responses is given.
	for {
		fmt.Print(prompt)
		reply, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		reply = strings.TrimSpace(strings.ToLower(reply))
		if reply == "" {
			reply = defaultEntry
		}

		for _, validResponse := range validResponses {
			if reply == validResponse {
				return reply, nil
			}
		}
	}
}

// promptListBool prompts the user for a boolean (yes/no) with the given prefix.
// The function will repeat the prompt to the user until they enter a valid
// response.
func promptListBool(reader *bufio.Reader, prefix string,
	defaultEntry string) (bool, error) { // nolint:unparam

	// Setup the valid responses.
	valid := []string{"n", "no", "y", "yes"}
	response, err := promptList(reader, prefix, valid, defaultEntry)
	if err != nil {
		return false, err
	}
	return response == "yes" || response == "y", nil
}

// promptBirhdayUnixTimeStamp prompts the user a Unix timestamp in second.
func promptUnixTimestamp(reader *bufio.Reader, prefix string,
	defaultEntry string) (time.Time, error) { // nolint:unparam

	prompt := fmt.Sprintf("%s [%s]: ", prefix, defaultEntry)

	for {
		fmt.Print(prompt)
		reply, err := reader.ReadString('\n')
		if err != nil {
			return time.Time{}, err
		}
		reply = strings.TrimSpace(strings.ToLower(reply))
		if reply == "" {
			reply = defaultEntry
		}
		ts, err := strconv.ParseInt(reply, 10, 64)
		if err != nil {
			fmt.Print(prompt)
			continue
		}

		return time.Unix(ts, 0), nil
	}
}

// promptPass prompts the user for a passphrase with the given prefix.  The
// function will ask the user to confirm the passphrase and will repeat the
// prompts until they enter a matching response.
func promptPass(_ *bufio.Reader, prefix string, confirm bool) ([]byte, error) {
	// Prompt the user until they enter a passphrase.
	prompt := fmt.Sprintf("%s: ", prefix)
	for {
		fmt.Print(prompt)
		pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return nil, err
		}
		fmt.Print("\n")
		pass = bytes.TrimSpace(pass)
		if len(pass) == 0 {
			continue
		}

		if !confirm {
			return pass, nil
		}

		fmt.Print("Confirm passphrase: ")
		confirm, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return nil, err
		}
		fmt.Print("\n")
		confirm = bytes.TrimSpace(confirm)
		if !bytes.Equal(pass, confirm) {
			fmt.Println("The entered passphrases do not match")
			continue
		}

		return pass, nil
	}
}

func birthday(reader *bufio.Reader) (time.Time, error) {
	prompt := "Enter the birthday of the seed in Unix timestamp " +
		"(the walllet will scan the chain from this time)"
	return promptUnixTimestamp(reader, prompt, "0")
}

// PrivatePass prompts the user for a private passphrase with varying behavior
// depending on whether the passed legacy keystore exists.  When it does, the
// user is prompted for the existing passphrase which is then used to unlock it.
// On the other hand, when the legacy keystore is nil, the user is prompted for
// a new private passphrase.  All prompts are repeated until the user enters a
// valid response.
func PrivatePass(reader *bufio.Reader, legacyKeyStore *keystore.Store) ([]byte, error) {
	// When there is not an existing legacy wallet, simply prompt the user
	// for a new private passphase and return it.
	if legacyKeyStore == nil {
		return promptPass(reader, "Enter the private "+
			"passphrase for your new wallet", true)
	}

	// At this point, there is an existing legacy wallet, so prompt the user
	// for the existing private passphrase and ensure it properly unlocks
	// the legacy wallet so all of the addresses can later be imported.
	fmt.Println("You have an existing legacy wallet.  All addresses from " +
		"your existing legacy wallet will be imported into the new " +
		"wallet format.")
	for {
		privPass, err := promptPass(reader, "Enter the private "+
			"passphrase for your existing wallet", false)
		if err != nil {
			return nil, err
		}

		// Keep prompting the user until the passphrase is correct.
		if err := legacyKeyStore.Unlock(privPass); err != nil {
			if err == keystore.ErrWrongPassphrase {
				fmt.Println(err)
				continue
			}

			return nil, err
		}

		return privPass, nil
	}
}

// PublicPass prompts the user whether they want to add an additional layer of
// encryption to the wallet.  When the user answers yes and there is already a
// public passphrase provided via the passed config, it prompts them whether or
// not to use that configured passphrase.  It will also detect when the same
// passphrase is used for the private and public passphrase and prompt the user
// if they are sure they want to use the same passphrase for both.  Finally, all
// prompts are repeated until the user enters a valid response.
func PublicPass(reader *bufio.Reader, privPass []byte,
	defaultPubPassphrase, configPubPassphrase []byte) ([]byte, error) {

	pubPass := defaultPubPassphrase
	usePubPass, err := promptListBool(reader, "Do you want "+
		"to add an additional layer of encryption for public "+
		"data?", "no")
	if err != nil {
		return nil, err
	}

	if !usePubPass {
		return pubPass, nil
	}

	if !bytes.Equal(configPubPassphrase, pubPass) {
		useExisting, err := promptListBool(reader, "Use the "+
			"existing configured public passphrase for encryption "+
			"of public data?", "no")
		if err != nil {
			return nil, err
		}

		if useExisting {
			return configPubPassphrase, nil
		}
	}
	pubPass, err = provideSeed(reader)
	if err != nil {
		return nil, err
	}

	fmt.Println("NOTE: Use the --walletpass option to configure your " +
		"public passphrase.")
	return pubPass, nil
}

// Seed prompts the user whether they want to use an existing wallet generation
// seed.  When the user answers no, a seed will be generated and displayed to
// the user along with prompting them for confirmation.  When the user answers
// yes, a the user is prompted for it.  All prompts are repeated until the user
// enters a valid response.
func Seed(reader *bufio.Reader) ([]byte, time.Time, error) {
	bday := time.Now()
	// Ascertain the wallet generation seed.
	useUserSeed, err := promptListBool(reader, "Do you have an "+
		"existing wallet seed you want to use?", "no")
	if err != nil {
		return nil, bday, err
	}
	if !useUserSeed {
		seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
		if err != nil {
			return nil, bday, err
		}

		fmt.Printf("Your wallet generation seed is: %x\n", seed)

		fmt.Printf("\nIMPORTANT: Keep the seed in a safe place as you " +
			"will NOT be able to restore your wallet without it.\n\n")

		for {
			fmt.Print(`Once you have stored the seed in a safe ` +
				`and secure location, enter "OK" to continue: `)
			confirmSeed, err := reader.ReadString('\n')
			if err != nil {
				return nil, bday, err
			}
			confirmSeed = strings.TrimSpace(confirmSeed)
			confirmSeed = strings.Trim(confirmSeed, `"`)
			if confirmSeed == "OK" {
				break
			}
		}

		return seed, bday, nil
	}

	seed, err := provideSeed(reader)
	if err != nil {
		return nil, bday, err
	}

	bday, err = birthday(reader)
	if err != nil {
		return nil, bday, err
	}

	return seed, bday, nil
}
