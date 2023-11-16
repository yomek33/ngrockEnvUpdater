package ngrockEnvUpdater

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// ngrok APIからのレスポンス
type NgrokTunnelsResponse struct {
	Tunnels []struct {
		PublicURL string `json:"public_url"`
	} `json:"tunnels"`
}

func EnvUpdate() {
	// ngrok APIからトンネル情報を取得
	resp, err := http.Get("http://127.0.0.1:4040/api/tunnels")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var tunnels NgrokTunnelsResponse
	err = json.Unmarshal(body, &tunnels)
	if err != nil {
		panic(err)
	}

	// .envファイルを読み込む
	envContent, err := ioutil.ReadFile(".env")
	if err != nil {
		panic(err)
	}

	// URLを更新
	newEnvContent := updateEnvFile(string(envContent), "NGROK_URL", tunnels.Tunnels[0].PublicURL)

	// .envファイルに書き込む
	err = ioutil.WriteFile(".env", []byte(newEnvContent), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Updated .env with new ngrok URL:", tunnels.Tunnels[0].PublicURL)
}

// .envファイルの内容を更新する
func updateEnvFile(content, key, newValue string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			lines[i] = key + "=" + newValue
			return strings.Join(lines, "\n")
		}
	}
	return content + "\n" + key + "=" + newValue
}
