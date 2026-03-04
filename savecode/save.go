
// var lockKeys []string
// 	uniqueKeys := make(map[string]bool)

// 	for _, item := range tx.Items {
// 		keyOrigin := fmt.Sprintf("lock:warehouse:%d:product:%d", tx.WarehouseID, item.ProductID)
// 		if !uniqueKeys[keyOrigin] {
// 			uniqueKeys[keyOrigin] = true
// 			lockKeys = append(lockKeys, keyOrigin)
// 		}

// 		if tx.Type == "TRANSFER" && tx.TargetWarehouseID != nil {
// 			keyTarget := fmt.Sprintf("lock:warehouse:%d:product:%d", *tx.TargetWarehouseID, item.ProductID)
// 			if !uniqueKeys[keyTarget] {
// 				uniqueKeys[keyTarget] = true
// 				lockKeys = append(lockKeys, keyTarget)
// 			}
// 		}
// 	}

// 	sort.Strings(lockKeys)

// 	var acquiredLocks []string
// 	defer func() {
// 		for _, key := range acquiredLocks {
// 			config.RedisClient.Del(config.Ctx, key)
// 		}
// 	}()

// 	for _, key := range lockKeys {
// 		isLocked, _ := config.RedisClient.SetNX(config.Ctx, key, "locked", 10*time.Second).Result()
// 		if !isLocked {
// 			return errors.New("Stok di salah satu gudang sedang diproses oleh sistem lain. Silakan coba lagi.")
// 		}
// 		acquiredLocks = append(acquiredLocks, key)
// 	}

// 	stockMutations := make(map[string]int)

// 	for _, item := range tx.Items {
// 		if tx.Type == "OUTBOUND" || tx.Type == "TRANSFER" {
// 			stockData, err := s.whStockRepo.GetStock(tx.WarehouseID, item.ProductID)
// 			if err != nil || stockData.Stock < item.Quantity {
// 				return fmt.Errorf("Persetujuan ditolak! Stok produk ID %d di Gudang Asal tidak mencukupi", item.ProductID)
// 			}

// 			originKey := fmt.Sprintf("%d_%d", tx.WarehouseID, item.ProductID)
// 			stockMutations[originKey] -= item.Quantity
// 		}

// 		if tx.Type == "INBOUND" {
// 			targetKey := fmt.Sprintf("%d_%d", tx.WarehouseID, item.ProductID)
// 			stockMutations[targetKey] += item.Quantity
// 		}

// 		if tx.Type == "TRANSFER" && tx.TargetWarehouseID != nil {
// 			targetKey := fmt.Sprintf("%d_%d", *tx.TargetWarehouseID, item.ProductID)
// 			stockMutations[targetKey] += item.Quantity
// 		}
// 	}

// 	tx.Status = "approved"
// 	tx.ApprovedByID = &approverID

// 	err = s.