# Teste de carga com k6


## 1. Recursos 

Restritos no `docker-compose.yml`:

* 0.5 CPU + 256MB RAM para aplicação web
* 0.5 CPU + 517MB RAM para banco de dados

## 2. Tags

- `versao_lenta`
- `versao_otimizada_1`: índice, paginação e otimizações de código
- `versao_otimizada_2`: pool de conexões

Cheque [o arquivo `app/amigos/handlers.go` na versão lenta](https://github.com/EdyKnopfler/teste-carga-k6/commit/97accfc9a872a461f62ca2b2129958ed5a812fa7) para ver problemas de performance e possíveis pontos de otimização documentados.

## 3. Scripts

- `limpar-e-popular.sql`: massa de dados inicial
- `rodar-stress-test.sh`: comando para rodar o k6

## 4. Etapas de otimização

Resultado inicial:

```
TOTAL RESULTS 


checks_total.......: 580     3.998776/s

checks_succeeded...: 100.00% 580 out of 580

checks_failed......: 0.00%   0 out of 580



✓ busca ok

✓ cadastro ok


HTTP

http_req_duration..............: avg=4.7s min=14.85ms med=2.45s max=29.87s p(90)=12.47s p(95)=17.23s

  { expected_response:true }...: avg=4.7s min=14.85ms med=2.45s max=29.87s p(90)=12.47s p(95)=17.23s

http_req_failed................: 0.00%  0 out of 580

http_reqs......................: 580    3.998776/s



EXECUTION

iteration_duration.............: avg=5.7s min=1.01s   med=3.45s max=30.88s p(90)=13.47s p(95)=18.23s

iterations.....................: 580    3.998776/s

vus............................: 1      min=1        max=50

vus_max........................: 50     min=50       max=50



NETWORK

data_received..................: 169 MB 1.2 MB/s

data_sent......................: 72 kB  494 B/s
```

### 4.1 Melhorar os algoritmos

* Tag: `versao_otimizada_1`

#### 4.1.1 O famigerado `ILIKE '%texto%'`

Usamos um [índice invertido](https://en.wikipedia.org/wiki/Inverted_index) (*GIN: Generalized Inverted Index*) para busca textual:

```sql
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_amigos_nome_trgm ON amigos USING gin (nome gin_trgm_ops);
```

#### 4.1.2 Paginação nas consultas

A primeira versão aceitava qualquer busca textual e devolvia todas as linhas encontradas.

#### 4.1.3 Evitando alocação dinâmica de memória

Quando já sabemos quantos elementos há no DTO, podemos alocar antecipadamente no objeto entidade (e vice-versa). Evitamos realocação dinâmica por `append` executado em loop, que num cenário de estresse pode forçar o Garbage Collector.

#### 4.1.4 Resultados

```
TOTAL RESULTS 


checks_total.......: 1693    12.052844/s

checks_succeeded...: 100.00% 1693 out of 1693

checks_failed......: 0.00%   0 out of 1693



✓ busca ok

✓ cadastro ok



HTTP

http_req_duration..............: avg=833.35ms min=3.64ms med=291.22ms max=7.7s p(90)=2.3s p(95)=3.83s

  { expected_response:true }...: avg=833.35ms min=3.64ms med=291.22ms max=7.7s p(90)=2.3s p(95)=3.83s

http_req_failed................: 0.00%  0 out of 1693

http_reqs......................: 1693   12.052844/s



EXECUTION

iteration_duration.............: avg=1.83s    min=1s     med=1.29s    max=8.7s p(90)=3.3s p(95)=4.83s

iterations.....................: 1693   12.052844/s

vus............................: 1      min=1         max=50

vus_max........................: 50     min=50        max=50



NETWORK

data_received..................: 3.2 MB 23 kB/s

data_sent......................: 204 kB 1.4 kB/s
```

### 4.2 Determinar tamanho ótimo do pool de conexões com banco de dados

* Tag: `versao_otimizada_2`

A falta de limites no pool de conexões causa contenção de recursos no PostgreSQL. Configuramos o tamanho ótimo do pool considerando o limite de hardware (0.5 vCPU):

* `SetMaxOpenConns(25)`: evita que o banco gaste mais tempo trocando contextos do que processando queries.
* `SetMaxIdleConns (10)`: mantém conexões prontas, eliminando o custo de handshake em cada requisição.


#### 4.2.1 Resultados

```
TOTAL RESULTS 


checks_total.......: 2286    16.235222/s

checks_succeeded...: 100.00% 2286 out of 2286

checks_failed......: 0.00%   0 out of 2286



✓ busca ok

✓ cadastro ok



HTTP

http_req_duration..............: avg=351.45ms min=2.92ms med=87.92ms max=3.7s p(90)=1.09s p(95)=1.9s

  { expected_response:true }...: avg=351.45ms min=2.92ms med=87.92ms max=3.7s p(90)=1.09s p(95)=1.9s

http_req_failed................: 0.00%  0 out of 2286

http_reqs......................: 2286   16.235222/s



EXECUTION

iteration_duration.............: avg=1.35s    min=1s     med=1.08s   max=4.7s p(90)=2.09s p(95)=2.9s

iterations.....................: 2286   16.235222/s

vus............................: 1      min=1         max=50

vus_max........................: 50     min=50        max=50



NETWORK

data_received..................: 4.3 MB 30 kB/s

data_sent......................: 281 kB 2.0 kB/s
```