package repository

import (
	"database/sql"
	"e-commerce-api/internal/errors"
	"e-commerce-api/internal/models"
	"fmt"
	"time"

	_"github.com/lib/pq"
)

type PostgresRepository struct{
	db *sql.DB
}

func NewPostgresRepository(connectionString string) (ProductRepository,error){
	db,err:=sql.Open("postgres",connectionString)
	if err!=nil{
		return nil,fmt.Errorf("failed to open database :%w",err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5*time.Minute)
	if err:=db.Ping();err!=nil{
		return nil,fmt.Errorf("failed to ping database :%w",err)
	}
	return &PostgresRepository{db: db},nil
}

func (r *PostgresRepository) Close()error{
	return r.db.Close()
}

func (r *PostgresRepository) GetAll()([]models.Product,error){
	query:="SELECT id,name,price from products ORDER BY id"
	rows,err:=r.db.Query(query)
	if err!=nil{
		return nil,fmt.Errorf("failed to query products :%w",err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next(){
		var product models.Product
		err:=rows.Scan(&product.ID,&product.Name,&product.Price)
		if err!=nil{
			return nil,fmt.Errorf("error while scanning the products: %w",err)
		}
		products=append(products, product)
	}
	if err =rows.Err();err!=nil{
		return nil,fmt.Errorf("row iteration errors:%w",err)
	}

	return products,nil
}

func (r *PostgresRepository) GetByID(id int)(*models.Product,error){
	query:=`select id,name,price from products where id =$1`
	var product models.Product
	err:=r.db.QueryRow(query,id).Scan(&product.ID,&product.Name,&product.Price)
	if err!=nil{
		if err==sql.ErrNoRows{
			return nil,errors.NewNotFoundError(id,"product")
		}
		return nil,fmt.Errorf("failed to get product :%w",err)
	}
	return &product,nil
}

func (r *PostgresRepository)Create(p *models.Product) error{
	query:=`insert into products (name,price) values ($1,$2) returning id`
	err:=r.db.QueryRow(query,p.Name,p.Price).Scan(&p.ID)
	if err!=nil{
		return fmt.Errorf("failed to create product :%w",err)

	}
	return nil
}
func (r *PostgresRepository)	Update(id int, p *models.Product) error{
	query :=` update products set name=$1, price=$2,updated_at=CURRENT_TIMESTAMP WHERE id=$3`
	result,err:=r.db.Exec(query,p.Name,p.Price,id)
	if err!=nil{
		return fmt.Errorf("failed to update the product :%w",err)
	}
	rowsAffected,err:=result.RowsAffected()
	if err!=nil{
		return fmt.Errorf("failed to get rows affected:%w",err)
	}
	if rowsAffected==0{
		return errors.NewNotFoundError(id,"product")
	}
	p.ID=id
	return nil
}
func (r *PostgresRepository)	Delete(id int) error{
	 query := `DELETE FROM products WHERE id = $1`
    
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete product: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return errors.NewNotFoundError(id,"product")
    }

    return nil
}
  func (r *PostgresRepository)  ExistsByName(name string)bool{
	query :=`select exists(select 1 from products where name=$1)`

	var exists bool
	err:=r.db.QueryRow(query,name).Scan(&exists)
	if err!=nil{
		return false
	}
	return exists
  }