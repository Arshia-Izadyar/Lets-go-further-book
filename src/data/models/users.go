package models

import "time"

type Users struct {
	ID        int64     
	CreatedAt time.Time 
	Name      string    
	Email     string    
	Activated bool      
	Version   int       
	Password  string  
}

var AnonymousUser = &Users{}
