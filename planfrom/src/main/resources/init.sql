-- agent_otc.agent definition

CREATE TABLE `agent` (
  `id` int NOT NULL AUTO_INCREMENT,
  `agent_name` varchar(255) NOT NULL,
  `wallet_address` varchar(255) NOT NULL,
  `encrypted_private_key` text NOT NULL,
  `token` varchar(255) NOT NULL,
  `status` varchar(50) NOT NULL DEFAULT 'ACTIVE',
  `profit_address` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `agent_name` (`agent_name`),
  UNIQUE KEY `token` (`token`),
  KEY `idx_agent_token` (`token`)
) ENGINE=InnoDB AUTO_INCREMENT=71 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- agent_otc.award_record definition

CREATE TABLE `award_record` (
  `id` int NOT NULL AUTO_INCREMENT,
  `award_id` varchar(255) NOT NULL,
  `agent_id` int NOT NULL,
  `award_type` varchar(100) NOT NULL,
  `description` text NOT NULL,
  `reward_amount` decimal(20,10) NOT NULL,
  `awarded_at` timestamp NOT NULL,
  `ranking` int NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `award_id` (`award_id`),
  KEY `idx_award_agent_id` (`agent_id`),
  CONSTRAINT `award_record_ibfk_1` FOREIGN KEY (`agent_id`) REFERENCES `agent` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- agent_otc.`order` definition

CREATE TABLE `order` (
  `id` int NOT NULL AUTO_INCREMENT,
  `order_id` varchar(255) NOT NULL,
  `agent_id` int NOT NULL,
  `type` varchar(50) NOT NULL,
  `token_address` varchar(255) NOT NULL,
  `token_symbol` varchar(50) NOT NULL,
  `amount` decimal(20,10) NOT NULL,
  `price` decimal(20,10) NOT NULL,
  `status` varchar(50) NOT NULL DEFAULT 'ACTIVE',
  `buyer_wallet` varchar(255) DEFAULT NULL,
  `seller_wallet` varchar(255) DEFAULT NULL,
  `transaction_hash` varchar(255) DEFAULT NULL,
  `chain_order_id` varchar(255) DEFAULT NULL,
  `error_msg` text,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `order_id` (`order_id`),
  KEY `idx_order_agent_id` (`agent_id`),
  KEY `idx_order_status` (`status`),
  CONSTRAINT `order_ibfk_1` FOREIGN KEY (`agent_id`) REFERENCES `agent` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=146 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


-- agent_otc.thought_process definition

CREATE TABLE `thought_process` (
  `id` int NOT NULL AUTO_INCREMENT,
  `thought_id` varchar(255) NOT NULL,
  `agent_id` int NOT NULL,
  `content` text NOT NULL,
  `type` varchar(50) NOT NULL,
  `timestamp` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `thought_id` (`thought_id`),
  KEY `idx_thought_agent_id` (`agent_id`),
  CONSTRAINT `thought_process_ibfk_1` FOREIGN KEY (`agent_id`) REFERENCES `agent` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=224 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


-- agent_otc.transaction_record definition

CREATE TABLE `transaction_record` (
  `id` int NOT NULL AUTO_INCREMENT,
  `transaction_id` varchar(255) NOT NULL,
  `agent_id` int NOT NULL,
  `type` varchar(50) NOT NULL,
  `token_symbol` varchar(50) NOT NULL,
  `amount` decimal(20,10) NOT NULL,
  `price` decimal(20,10) NOT NULL,
  `total_value` decimal(20,10) NOT NULL,
  `status` varchar(50) NOT NULL,
  `timestamp` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `transaction_id` (`transaction_id`),
  KEY `idx_transaction_agent_id` (`agent_id`),
  CONSTRAINT `transaction_record_ibfk_1` FOREIGN KEY (`agent_id`) REFERENCES `agent` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;