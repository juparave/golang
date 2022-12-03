# Go HTTP

https://www.digitalocean.com/community/tutorials/how-to-make-http-requests-in-go

## Post multipart/form-data using http.NewRequest

ref: https://golangbyexample.com/http-mutipart-form-body-golang/

```go
// HTTPPostForm posts Form payload to url using multipart/form-data
func HTTPPostForm(myURL string, payload map[string]string) (string, error) {

	// new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range payload {
		fw, err := writer.CreateFormField(k)
		if err != nil {
			// what!!
		}
		_, err = io.Copy(fw, strings.NewReader(string(v)))
		if err != nil {
			panic(err)
		}

	}
	// Close multipart writer.
	writer.Close()

	fmt.Println(body)

	request, err := http.NewRequest(http.MethodPost, myURL, bytes.NewReader(body.Bytes()))
	// request.Header.Set("Content-Type", "multipart/form-data; charset=UTF-8")
	// use `writer` generated header content, it will add boundary value
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		response.Body.Close()
		return "", err
	}
	defer response.Body.Close()

	resBody, _ := ioutil.ReadAll(response.Body)

	return string(resBody), nil
}


```
