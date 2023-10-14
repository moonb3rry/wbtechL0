package pgrepo

import (
	"WBTech0/internal/entity"
	"WBTech0/pkg/postgres"
	"context"
	"fmt"
)

type OrderRepo struct {
	db *postgres.Postgres
}

func NewOrderRepo(db *postgres.Postgres) *OrderRepo {
	return &OrderRepo{db}
}

func (r OrderRepo) GetAllOrders(ctx context.Context) ([]entity.Order, error) {
	query := `SELECT 
        o.order_uid,
        o.track_number,
        o.entry,
        (SELECT d FROM delivery d WHERE d.id = o.delivery_id) AS delivery,
        (SELECT p FROM payment p WHERE p.transaction = o.payment_id) AS payment,
        (SELECT i 
         FROM item i
         JOIN items_in_order io ON i.chrt_id = io.item_id
         WHERE io.order_id = o.order_uid) AS items,
        o.locale,
        o.internal_signature,
        o.customer_id,
        o.delivery_service,
        o.shardkey,
        o.sm_id,
        o.date_created,
        o.oof_shard
    FROM orders o`

	rows, err := r.db.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	var ordersMap = make(map[string]entity.Order)
	for rows.Next() {
		var order entity.Order
		var item entity.Item
		var orderUID string
		if err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID,
			&order.DeliveryService, &order.ShardKey, &order.SMID, &order.DateCreated, &order.OOFShard,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
			&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDT, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee,
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NMID, &item.Brand, &item.Status); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		if existingOrder, ok := ordersMap[orderUID]; ok {
			existingOrder.Items = append(existingOrder.Items, item)
			ordersMap[orderUID] = existingOrder
		} else {
			order.Items = append(order.Items, item)
			ordersMap[orderUID] = order
		}
	}
	var orders []entity.Order
	for _, order := range ordersMap {
		orders = append(orders, order)
	}

	return orders, nil

}

//func (r *OrderRepo) GetOrderById(ctx context.Context, orderUID string) (entity.Order, error) {
//	order := entity.Order{}
//	item := entity.Item{}
//	query := `SELECT
//        o.order_uid,
//        o.track_number,
//        o.entry,
//        (SELECT d FROM delivery d WHERE d.id = o.delivery_id) AS delivery,
//        (SELECT p FROM payment p WHERE p.transaction = o.payment_id) AS payment,
//        (SELECT i
//         FROM item i
//         JOIN items_in_order io ON i.chrt_id = io.item_id
//         WHERE io.order_id = o.order_uid) AS items,
//        o.locale,
//        o.internal_signature,
//        o.customer_id,
//        o.delivery_service,
//        o.shardkey,
//        o.sm_id,
//        o.date_created,
//        o.oof_shard
//    FROM orders o
//    WHERE o.order_uid = $1`
//
//	row := r.db.Pool.QueryRow(context.Background(), query, orderUID)
//
//	if err := row.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID,
//		&order.DeliveryService, &order.ShardKey, &order.SMID, &order.DateCreated, &order.OOFShard,
//		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
//		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDT, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee,
//		&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NMID, &item.Brand, &item.Status,
//	); err != nil {
//		return entity.Order{}, fmt.Errorf("failed to scan row: %v", err)
//	}
//	return order, nil
//}

func (r *OrderRepo) AddOrder(ctx context.Context, order entity.Order) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// delivery
	var deliveryID int
	err = tx.QueryRow(ctx, `INSERT INTO delivery(name, phone, zip, city, address, region, email) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email).Scan(&deliveryID)
	if err != nil {
		return err
	}

	// payment
	var paymentTransaction string
	err = tx.QueryRow(ctx, `INSERT INTO payment(transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING transaction`,
		order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee).Scan(&paymentTransaction)
	if err != nil {
		return err
	}

	// items
	var itemIDs []int
	for _, item := range order.Items {
		var itemID int
		err = tx.QueryRow(ctx, `INSERT INTO item(chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING chrt_id`,
			item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NMID, item.Brand, item.Status).Scan(&itemID)
		if err != nil {
			return err
		}
		itemIDs = append(itemIDs, itemID)
	}

	// order
	_, err = tx.Exec(ctx, `INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, delivery_id, payment_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SMID, order.DateCreated, order.OOFShard, deliveryID, paymentTransaction)
	if err != nil {
		return err
	}

	// items in order
	for _, itemID := range itemIDs {
		_, err = tx.Exec(ctx, `INSERT INTO items_in_order(order_id, item_id) VALUES($1, $2)`, order.OrderUID, itemID)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
