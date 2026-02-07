# Teste de carga com k6

## 1. Tags

- `versao_lenta`
- `versao_otimizada_xxx` **(TODO)**

Cheque [o arquivo `app/amigos/handlers.go` na versão lenta](https://github.com/EdyKnopfler/teste-carga-k6/commit/97accfc9a872a461f62ca2b2129958ed5a812fa7) para ver problemas de performance e possíveis pontos de otimização documentados.

## 2. Scripts

- `limpar-e-popular.sql`: massa de dados inicial
- `rodar-stress-test.sh`: comando para rodar o k6

## 3. Etapas de otimização

### 3.1 Melhorar os algoritmos

#### 3.1.1 O famigerado `ILIKE '%texto%'`

Usamos um [índice invertido](https://en.wikipedia.org/wiki/Inverted_index) (*GIN: Generalized Inverted Index*) para busca textual:

```sql
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_amigos_nome_trgm ON amigos USING gin (nome gin_trgm_ops);
```

# 3.1.2 Paginação nas consultas

A primeira versão aceitava qualquer busca textual e devolvia todas as linhas encontradas.

### 3.2 Determinar tamanho ótimo do pool