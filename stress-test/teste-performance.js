import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 20 }, // Sobe para 20 usuários
    { duration: '1m', target: 20 },  // Mantém 20 usuários (estável)
    { duration: '30s', target: 50 }, // Stress: sobe para 50
    { duration: '20s', target: 0 },  // Recuperação
  ],
};

const BASE_URL = 'http://host.docker.internal:8080';

export default function() {
  // 1. Simula Cadastro (Escrita) - 20% das vezes
  if (Math.random() < 0.2) {
    let payload = JSON.stringify({
      nome: `Amigo ${__VU}-${__ITER}`,
      dataNascimento: '1990-01-01',
      preferencias: ['Go', 'Docker', 'Performance']
    });
    
    let params = { headers: { 'Content-Type': 'application/json' } };
    let res = http.post(`${BASE_URL}/amigos`, payload, params);
    check(res, { 'cadastro ok': (r) => r.status === 201 });
  }

  // 2. Simula Busca (Leitura + ILIKE) - 80% das vezes
  else {
    // Termos aleatórios para evitar que o banco cacheie tudo
    const termos = ['Amigo', 'Performance', 'a', 'z']; 
    const termo = termos[Math.floor(Math.random() * termos.length)];
    
    let res = http.get(`${BASE_URL}/amigos?q=${termo}`);
    check(res, { 'busca ok': (r) => r.status === 200 });
  }

  sleep(1); // Simula o "think time" do usuário
}