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
	"golang.org/x/crypto/ssh/terminal"
)

func promptSeed(reader *bufio.Reader) ([]byte, error) {
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

// promptPassphrase prompts the user for a passphrase with the given prefix.
// The function will ask the user to confirm the passphrase and will repeat
// the prompts until they enter a matching response.
func promptPassphrase(prefix string, confirm bool) ([]byte, error) {
	pass := os.Getenv("LBCWALLET_PASSPHRASE")
	if len(pass) > 0 {
		return []byte(pass), nil
	}

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

// Passphrase prompts the user for a passphrase.
// All prompts are repeated until the user enters a valid response.
func Passphrase(confirm bool) ([]byte, error) {
	return promptPassphrase("Enter the passphrase "+
		"for your new wallet", confirm)
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

	seed, err := promptSeed(reader)
	if err != nil {
		return nil, bday, err
	}

	bday, err = birthday(reader)
	if err != nil {
		return nil, bday, err
	}

	return seed, bday, nil
}
