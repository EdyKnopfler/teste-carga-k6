# Teste de carga com k6

## Tags

- `versao_lenta`
- `versao_otimizada_xxx` **(TODO)**

Cheque [o arquivo `app/amigos/handlers.go` na versão lenta](https://github.com/EdyKnopfler/teste-carga-k6/commit/97accfc9a872a461f62ca2b2129958ed5a812fa7) para ver problemas de performance e possíveis pontos de otimização documentados.

## Scripts

- `limpar-e-popular.sql`: massa de dados inicial
- `rodar-stress-test.sh`: comando para rodar o k6

## Etapas de otimização

* [ ] Determinar tamanho ótimo do pool
* [ ] Melhorar os algoritmos