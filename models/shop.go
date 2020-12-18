package model

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

type Shop struct {
	Model

	/** shop id */
	//ShopId string `gorm:"type:varchar(36);column:shop_id;" json:"shop_id" form:"shop_id"`

	/** shop title */
	Title string `gorm:"type:varchar(255);" form:"title" binding:"required" json:"title"`

	/** 商户id merchant表的merchant_id */
	MerchantId uint64 `gorm:"type:bigint(20);column:merchant_id;" form:"merchantId" binding:"required" json:"merchantId"`

	/** 店铺唯一标识 */
	Code string `gorm:"type:varchar(20);" form:"code" json:"-"`

	/** 评分*/
	Star int8 `gorm:"type:int(10);" form:"star" json:"star"`

	/** 营业时间*/
	BusinessTime string `gorm:"type:varchar(100);" form:"businessTime" json:"businessTime"`

	/** 关注人数*/
	AttentionNum uint64 `gorm:"type:bigint(20);column:attention_num;" form:"attentionNum" json:"attentionNum"`

	/** 店铺电话 */
	Phone string `gorm:"type:varchar(20);" form:"phone" binding:"required" json:"phone"`

	/** shop address */
	Address string `gorm:"type:varchar(255);" form:"address" json:"address"`

	/** 店铺地址描述 */
	AddressDescription string `gorm:"type:varchar(255);" form:"addressDescription" json:"addressDescription"`

	/** 店铺描述 */
	Description string `gorm:"type:varchar(1000);" form:"description" json:"description"`

	/** 店铺图文详情 */
	DetailDescribe string `gorm:"type:varchar(5000);" form:"detailDescribe" json:"detail_describe"`

	/** 店铺logo */
	Logo string `gorm:"type:varchar(500);" form:"logo" json:"logo"`

	/** shop url */
	ImageUrl string `gorm:"type:varchar(2000);" form:"image_url" json:"image_url"`

	/** shop 的描述 */
	Content string `gorm:"type:varchar(255)" form:"content" json:"content"`

	StringDuration

	Longitude float64 `gorm:"type:decimal" form:"longitude" json:"longitude"`
	Latitude  float64 `gorm:"type:decimal" form:"latitude" json:"latitude"`

	ProductList []*Product `gorm:"ForeignKey:ShopId;association_foreignkey:ShopId" json:"product_list"`

	Coordinates GeoPoint

	/** 0:店铺未开通，1:店铺可用，2:冻结，3：永久封店，4：其他*/
	Status uint8 `gorm:"type:tinyint(1);default:1" form:"status" json:"status"`

	/**店铺类型 1:平台官方店，2:商户代平台核销的店， 3：商户官方店，4：商户直营店，5：其他品牌店*/
	Type uint8 `gorm:"type:tinyint(1);default:3" form:"type" json:"type"`

	/** 越大越往前排*/
	OrderWeight uint64 `gorm:"type:bigint;default:1" form:"order_weight" json:"order_weight"`

	//Location GeoPoint `gorm:"column:location;type:geometry; not null" sql:"type:geometry(Geometry,4326)"`

	CrudTime
}

type ShopDetail struct {
	/** shop id */
	ShopId uint64 `gorm:"type:bigint(20);" json:"shopId"`
	Title  string `gorm:"type:varchar(255);" json:"title"`
	Logo   string `gorm:"type:varchar(500);" json:"logo"`
}

type ShopDetailDistance struct {
	ShopDetail
	Longitude float64 `gorm:"column:longitude;" form:"longitude" json:"longitude"`
	Latitude  float64 `gorm:"column:latitude;" form:"latitude" json:"latitude"`
	Distance  float64 `gorm:"column:distance;" json:"distance"`
}

// 表结构初始化
func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&Shop{})
}

func (ShopDetail) TableName() string {
	//return constant.ShopDetailTable
	tableName := helper.GetDBByName(constant.DBMerchant).Model(&ShopDetail{}).Name()
	return tableName
}

// 插入前生成主键
func (shop *Shop) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//shop.ID = id.String()
	return nil
}

// 校验表单中提交的参数是否合法
func (shop *Shop) Validator() error {
	return nil
}

type GeoPoint struct {
	Longitude float64 `form:"longitude" json:"longitude"`
	Latitude  float64 `form:"latitude" json:"latitude"`
}

func (p *GeoPoint) String2() string {
	point := fmt.Sprintf("SRID=4326;POINT(%v %v)", p.Longitude, p.Latitude)
	return fmt.Sprintf("GeomFromText('%v')", point)
}

func (p *GeoPoint) String() string {
	return fmt.Sprintf("POINT(%v %v)", p.Longitude, p.Latitude)
}

// Scan implements the Scanner interface which will scan the postgis POINT(x, y) into the GeoPoint struct
func (p *GeoPoint) Scan(val interface{}) error {
	b, err := hex.DecodeString(string(val.([]uint8)))
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)
	var wkbByteOrder uint8
	if err := binary.Read(r, binary.LittleEndian, &wkbByteOrder); err != nil {
		return err
	}

	var byteOrder binary.ByteOrder
	switch wkbByteOrder {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return fmt.Errorf("invalid byte order %v", wkbByteOrder)
	}

	var wkbGeometryType uint64
	if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
		return err
	}

	if err := binary.Read(r, byteOrder, p); err != nil {
		return err
	}

	return nil
}

func (p GeoPoint) Value() (driver.Value, error) {
	return p.String(), nil
}
