-- Limpeza
TRUNCATE TABLE amigos, preferencia RESTART IDENTITY CASCADE;

-- Insere 10.000 amigos com nomes variados para o ILIKE
INSERT INTO amigos (id, nome, data_nascimento)
SELECT 
    gen_random_uuid(), 
    (ARRAY['Ana', 'Bruno', 'Carlos', 'Daniela', 'Eduardo', 'Fernanda', 'Gabriel', 'Helena'])[floor(random() * 8 + 1)] || ' ' || 
    (ARRAY['Silva', 'Santos', 'Oliveira', 'Souza', 'Pereira', 'Ferreira', 'Almeida', 'Costa'])[floor(random() * 8 + 1)] || ' ' || s.i,
    '1980-01-01'::date + (random() * 15000 * interval '1 day')
FROM generate_series(1, 10000) s(i);

-- Insere 2 a 3 preferências aleatórias para cada amigo
INSERT INTO preferencia (id, id_amigo, nome)
SELECT
    gen_random_uuid(),
    amigos.id, 
    (ARRAY['Go', 'Java', 'Python', 'Docker', 'Kubernetes', 'Cloud', 'Rust', 'SQL'])[floor(random() * 8 + 1)]
FROM amigos, generate_series(1, (random() * 2 + 1)::int);