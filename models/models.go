package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID `json:"_id" bson:"_id"`
	First_Name       *string            `json:"first_name" bson:"First_Name"`
	Last_Name        *string            `json:"last_name" bson:"Last_Name"`
	Password         *string            `json:"password" bson:"Password"`
	Email            *string            `json:"email" bson:"Email"`
	Phone            *int64             `json:"phone" bson:"Phone"`
	Token            *string            `json:"token" bson:"Token"`
	Refresh_Token    *string            `json:"refresh_token" bson:"Refresh_Token"`
	Created_At       time.Time          `json:"created_at" bson:"Created_At"`
	Updated_At       time.Time          `json:"updated_at" bson:"Updated_At"`
	User_ID          *string            `json:"user_id" bson:"user_id"`
	UserCart         []ProductUser      `json:"usercart" bson:"UserCart"`
	Address_Details  []Address          `json:"address_details" bson:"Address_Details"`
	Order_Status     []Order            `json:"order_status" bson:"Order_Status"`
}

type Product struct {
	Product_ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Product_Name *string            `json:"product_name" bson:"Product_Name"`
	Price        *uint64            `json:"price" bson:"Price"`
	Rating       *uint64            `json:"rating" bson:"Rating"`
	Image        *string            `json:"image" bson:"Image"`
}

type ProductUser struct {
	Product_ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Product_Name *string            `json:"product_name" bson:"Product_Name"`
	Price        *uint64            `json:"price" bson:"Price"`
	Rating       *uint64            `json:"rating" bson:"Rating"`
	Image        *string            `json:"image" bson:"Image"`
}

type Address struct {
	Address_ID primitive.ObjectID `json:"_id" bson:"_id"`
	House      *string            `json:"house" bson:"House"`
	Street     *string            `json:"street" bson:"Street"`
	City       *string            `json:"city" bson:"City"`
	Pincode    int64              `json:"pincode" bson:"PinCode"`
}

type Order struct {
	Order_ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Order_Cart      []ProductUser      `json:"order_cart" bson:"Order_Cart"`
	Order_At        time.Time          `json:"order_at" bson:"Order_At"`
	Price           int                `json:"price" bson:"Price"`
	Discount        int                `json:"discount" bson:"Discount"`
	Payment_Method  Payment            `json:"payment_method" bson:"Payment_Method"`
}

type Payment struct {
	Digital bool `json:"digital" bson:"Digital"`
	COD     bool `json:"cod" bson:"COD"`
}
