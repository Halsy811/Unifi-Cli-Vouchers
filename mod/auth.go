package auth

import (
	"fmt"
	"os"
	"syscall"

	"github.com/zalando/go-keyring" // для сохранения пароля в системном хранилище
	"golang.org/x/term"
)