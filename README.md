# Comparei - Serviço Proconf (Processamento e Confirmação) ⚙️

O **Comparei - Serviço Proconf** é um microsserviço do tipo *worker* dedicado a tarefas de processamento em *background* dentro do ecossistema "Comparei". 

Diferente de uma API tradicional, este serviço não expõe rotas HTTP. A sua função é atuar nos bastidores, ouvindo eventos de mensageria em tempo real e executando rotinas agendadas (CRON jobs), como o cálculo automático da pontuação de confiança (*Confidence Scores*) dos produtos e processamento de confirmações de utilizadores.

## 🛠️ Tecnologias Utilizadas

* **Linguagem:** [Go 1.23+](https://golang.org/)
* **Base de Dados NoSQL:** **MongoDB** (para armazenamento flexível de coleções como `usuarios`, `logs_confirmacao` e `mercado_produtos`).
* **Mensageria em Cache:** **Redis** (para arquitetura orientada a eventos).
* **Agendamento de Tarefas:** [Robfig Cron v3](https://github.com/robfig/cron) (para execução de *jobs* a cada X horas).
* **Infraestrutura:** Docker e Docker Compose.

## ⚙️ Arquitetura do Sistema

A aplicação inicializa dois motores principais:
1. **Subscriber (Eventos):** Ouve canais no Redis para processar instantaneamente ações de utilizadores, logs e configurações de produtos.
2. **Cron Scheduler:** Um gestor de tarefas que corre em segundo plano. Por exemplo, a cada 6 horas (`0 */6 * * *`), ele executa a rotina `CalculateConfidenceScores()`.

## 🚀 Como Executar o Projeto Localmente

### Pré-requisitos
* [Go 1.23+](https://golang.org/dl/) instalado no teu ambiente local.
* [Docker](https://www.docker.com/) e [Docker Compose](https://docs.docker.com/compose/) instalados para iniciar os serviços de base de dados e mensageria.

### Passo a Passo

1. **Clonar o repositório:**
```bash
   git clone https://github.com/ViniciusPereira03/comparei-servico-proconf
   cd comparei-servico-proconf

```

2. **Configuração das Variáveis de Ambiente:**
Cria um ficheiro `.env` na raiz do projeto. **Atenção:** Embora possas ter um `env.example` antigo, a aplicação precisa das variáveis do MongoDB e Redis para funcionar corretamente. Preenche da seguinte forma para correr localmente:
```env
# Configurações do MongoDB
MONGO_URI=mongodb://root:proconfdb@localhost:27017
MONGO_DB_NAME=proconfdb

# Configurações do Redis
REDIS_MESSAGING_HOST=localhost
REDIS_MESSAGING_PORT=6379

# Opcional (se existirem portas no docker-compose)
PORT=8085

```

3. **Executar a Aplicação com `run.sh`:**
Dá a permissão de execução ao script e inicializa o projeto:
```bash
chmod +x run.sh
./run.sh

```

> **⚠️ Nota Importante sobre o `wait-for-it.sh`:** > O script de arranque está configurado para utilizar a ferramenta `wait-for-it.sh`. A principal função deste script é colocar a aplicação em compasso de espera até que o MongoDB e o Redis estejam 100% inicializados e prontos para aceitar ligações. Isto garante que a aplicação não crashe logo no arranque com erros de "ligação recusada" (connection refused).


4. **Acompanhar os Logs:**
Após a execução, deves ver mensagens no terminal indicando:
* `📡 Inicializando subscriber...`
* `🚀 Serviço Proconf iniciado e aguardando eventos/cron...`
* Mensagens do CRON indicando recursos de memória e *Goroutines* ativas.



## 📂 Estrutura de Diretórios (Resumo)

* `/config`: Carregamento do ficheiro `.env`.
* `/internal`: Lógica central da aplicação.
    * `/app`: Serviços de domínio (`proconf_service.go`, `logs_service.go`, etc).
    * `/domain`: Entidades e interfaces.
    * `/infrastructure`:
        * `/messaging`: Subscritores que ouvem o Redis (`subscriber/proconf.go`, etc).
        * `/repository`: Operações diretas com o MongoDB (`proconf_repo.go`).
