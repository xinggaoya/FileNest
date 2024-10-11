package model

import (
	"gorm.io/gorm"
	"time"
)

/**
  @author: XingGao
  @date: 2024/9/28
**/

type BaseEntity struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}
