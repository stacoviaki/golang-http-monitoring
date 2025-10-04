## 🧭 COMO USAR
Para executar o monitoramento, utilize o comando abaixo no terminal:
```bash
go run main.go <server_list.csv> <downtime.csv>
```

Observações importantes: os dois parâmetros são obrigatórios. O primeiro arquivo (server_list.csv) deve existir e conter a lista de sites a serem verificados. O segundo arquivo (downtime.csv) será criado automaticamente com os sites que apresentarem falha.

📂 Parâmetros: server_list.csv — arquivo com a lista de sites que serão scaneados. Cada linha deve conter apenas uma URL, por exemplo:
NomeDoServidor,http://www.nomedoservidor.com/
