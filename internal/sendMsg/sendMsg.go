package sendMsg

import (
	"bot/internal/utils"
	"bot/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

//ctx这部分未做好

func SendMsg(pool *models.WorkPool) func(c *gin.Context) {
	return func(c *gin.Context) {
		event, _ := utils.GetEvent(c)
		go pool.AddEvent(&event)
		c.JSON(http.StatusOK, gin.H{})
	}
}
