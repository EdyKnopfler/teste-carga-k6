package amigos

// *** SÓ MANJA QUANTOS PROBLEMAS DE PERFORMANCE, POTENCIAIS OU REAIS, TEMOS NESTE ARQUIVO!

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type cadastroAmigos struct {
	DB *gorm.DB
}

type AmigoDTO struct {
	ID             string   `json:"id,omitempty"`
	Nome           string   `json:"nome" binding:"required"`
	DataNascimento string   `json:"dataNascimento" binding:"required"`
	Preferencias   []string `json:"preferencias"`
}

func Inicializar(DB *gorm.DB, router *gin.Engine) {
	amigos := &cadastroAmigos{DB: DB}

	DB.AutoMigrate(&Amigo{})
	DB.AutoMigrate(&Preferencia{})

	router.POST("/amigos", amigos.inserir)
	router.GET("/amigos", amigos.buscar)
	router.GET("/amigos/:id", amigos.buscarPorId)
}

func (cadastro *cadastroAmigos) inserir(c *gin.Context) {
	var dto AmigoDTO

	// Ficar (des)serializando tudo em JSON vai consumir boa parte da CPU
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	dataNasc, err := time.Parse("2006-01-02", dto.DataNascimento)

	if err != nil {
		c.JSON(400, gin.H{"error": "Formato de data inválido. Use AAAA-MM-DD"})
		return
	}

	novoAmigo := Amigo{
		Nome:           dto.Nome,
		DataNascimento: dataNasc,
	}

	// Já sei quantas preferências há! Posso pré-alocar.
	for _, p := range dto.Preferencias {
		novoAmigo.Preferencias = append(novoAmigo.Preferencias, Preferencia{
			Nome: p,
		})
	}

	// Peso do ORM
	if err := cadastro.DB.Create(&novoAmigo).Error; err != nil {
		c.JSON(500, gin.H{"error": "Erro ao salvar no banco: " + err.Error()})
		return
	}

	c.JSON(201, novoAmigo)
}

func (cadastro *cadastroAmigos) buscar(c *gin.Context) {
	termo := c.Query("q")

	if termo == "" {
		c.JSON(400, gin.H{"error": "O parâmetro de busca 'q' é obrigatório"})
		return
	}

	var amigos []Amigo

	// A parte que vai ficar "feia" no teste de carga
	searchQuery := "%" + termo + "%"

	// Preload pode gerar um IDAmigo IN (...) bem longo
	// Nada de paginação
	if err := cadastro.DB.Preload("Preferencias").Where("nome ILIKE ?", searchQuery).Find(&amigos).Error; err != nil {
		c.JSON(500, gin.H{"error": "Erro ao realizar busca: " + err.Error()})
		return
	}

	// Outra coisa que pode atrapalhar:
	amigosDTO := make([]AmigoDTO, 0)

	for _, amigo := range amigos {
		amigosDTO = append(amigosDTO, mapToDTO(amigo))
	}

	/*
		// Mais eficiente: evita realocações, GC, ...
		amigosDTO := make([]AmigoDTO, len(amigos))

		for i := range amigos {
			amigosDTO[i] = mapToDTO(amigos[i])
		}
	*/

	c.JSON(200, amigosDTO)
}

func (cadastro *cadastroAmigos) buscarPorId(c *gin.Context) {
	id := c.Param("id")

	var amigo Amigo

	if err := cadastro.DB.Preload("Preferencias").First(&amigo, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Amigo não encontrado"})
			return
		}

		c.JSON(500, gin.H{"error": "Erro ao buscar no banco: " + err.Error()})
		return
	}

	c.JSON(200, mapToDTO(amigo))
}

func mapToDTO(amigo Amigo) AmigoDTO {
	preferencias := make([]string, len(amigo.Preferencias))

	for i, p := range amigo.Preferencias {
		preferencias[i] = p.Nome
	}

	return AmigoDTO{
		ID:             amigo.ID.String(),
		Nome:           amigo.Nome,
		DataNascimento: amigo.DataNascimento.Format("2006-01-02"),
		Preferencias:   preferencias,
	}
}
