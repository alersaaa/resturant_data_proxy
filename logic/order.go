package logic

import (
	"fmt"
	"source/client"
	"strconv"
	"time"

	pb "source/pb"
)

// AddOrder 下订单
func AddOrder(req *pb.AddOrderReq, rsp *pb.AddOrderRsp) error {
	// 总点单数
	val, err := client.Rdb.Get(req.User.Username).Result()
	if err != nil {
		fmt.Printf("get count failed, err:%v\n", err)
		return err
	}
	count, err := strconv.Atoi(val)
	if err != nil {
		fmt.Printf("get count failed, err:%v\n", err)
		return err
	}
	expire_time := time.Now().AddDate(0, 0, 7)
	for _, n := range req.Products {
		// 各商品被点单数
		val1, err := client.Rdb.Get(req.User.Username + n.ProductName).Result()
		if err != nil {
			fmt.Printf("get product count failed, err:%v\n", err)
			client.Rdb.Set(req.User.Username+n.ProductName, n.Nums, time.Until(expire_time))
		}
		count1, _ := strconv.Atoi(val1)
		client.Rdb.Set(req.User.Username+n.ProductName, count1+int(n.Nums), time.Until(expire_time))
		count += int(n.Nums)
	}
	err = client.Rdb.Set(req.User.Username, count+1, time.Until(expire_time)).Err()
	if err != nil {
		fmt.Printf("set count failed, err:%v\n", err)
		return err
	}

	return nil
}
