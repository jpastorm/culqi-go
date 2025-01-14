package culqi

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	apiVersion    = "v2.0"
	baseURL       = "https://api.culqi.com/v2"
	baseSecureURL = "https://secure.culqi.com/v2"
)

// Errors API
var (
	ErrInvalidRequest = errors.New("La petición tiene una sintaxis inválida")
	ErrAuthentication = errors.New("La petición no pudo ser procesada debido a problemas con las llaves")
	ErrParameter      = errors.New("Algún parámetro de la petición es inválido")
	ErrCard           = errors.New("No se pudo realizar el cargo a una tarjeta")
	ErrLimitAPI       = errors.New("Estás haciendo muchas peticiones rápidamente al API o superaste tu límite designado")
	ErrResource       = errors.New("El recurso no puede ser encontrado, es inválido o tiene un estado diferente al permitido")
	ErrAPI            = errors.New("Error interno del servidor de Culqi")
	ErrUnexpected     = errors.New("Error inesperado, el código de respuesta no se encuentra controlado")
)

// WrapperResponse respuesta generica para respuestas GetAll
type WrapperResponse struct {
	Paging struct {
		Previous string `json:"previous"`
		Next     string `json:"next"`
		Cursors  struct {
			Before string `json:"before"`
			After  string `json:"after"`
		} `json:"cursors"`
	} `json:"paging"`
}

func do(method, endpoint string, params url.Values, body io.Reader) ([]byte, error) {
	if len(params) != 0 {
		endpoint += "?" + params.Encode()
	}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+keyInstance.SecretKey)

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	obj, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case 400:
		err = ErrInvalidRequest
	case 401:
		err = ErrAuthentication
	case 422:
		err = ErrParameter
	case 402:
		err = ErrCard
	case 429:
		err = ErrLimitAPI
	case 404:
		err = ErrResource
	case 500, 503:
		err = ErrAPI
	}

	if err != nil {
		err = fmt.Errorf("%v: %s", err, string(obj))
		return nil, err
	}

	if res.StatusCode >= 200 && res.StatusCode <= 206 {
		return obj, nil
	}

	return nil, ErrUnexpected
}

func doSecure(method, endpoint string, params url.Values, body io.Reader) ([]byte, error) {
	if len(params) != 0 {
		endpoint += "?" + params.Encode()
	}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+keyInstance.PublicKey)

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	obj, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case 400:
		err = ErrInvalidRequest
	case 401:
		err = ErrAuthentication
	case 422:
		err = ErrParameter
	case 402:
		err = ErrCard
	case 429:
		err = ErrLimitAPI
	case 404:
		err = ErrResource
	case 500, 503:
		err = ErrAPI
	}

	if err != nil {
		err = fmt.Errorf("%v: %s", err, string(obj))
		return nil, err
	}

	if res.StatusCode >= 200 && res.StatusCode <= 206 {
		return obj, nil
	}

	return nil, ErrUnexpected
}
