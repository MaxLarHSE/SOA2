CREATE TYPE product_status AS ENUM ('ACTIVE', 'INACTIVE', 'ARCHIVED');

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- При любом UPDATE меняем время на текущее
RETURN NEW;             -- Возвращаем обновленную строку
END;
$$ language 'plpgsql';

-- 3. Создаем саму таблицу products
CREATE TABLE products (
    -- gen_random_uuid() автоматически генерирует UUID при INSERT
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      name VARCHAR(255) NOT NULL,
      description VARCHAR(4000),
    -- DECIMAL(12,2) означает: максимум 12 цифр всего, из них 2 после запятой (например, 9999999999.99)
      price DECIMAL(12,2) NOT NULL,
      stock INTEGER NOT NULL DEFAULT 0,
      category VARCHAR(100) NOT NULL,
    -- Используем наш ENUM, по умолчанию товар будет 'ACTIVE'
      status product_status NOT NULL DEFAULT 'ACTIVE',
    -- seller_id пока может быть NULL (так как пользователей у нас еще нет)
      seller_id UUID,
    -- При создании записи автоматически ставим текущее время
      created_at TIMESTAMP NOT NULL DEFAULT NOW(),
      updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 4. Создаем индекс на поле status (требование из пункта 4)
-- Это сильно ускорит запросы вида: SELECT * FROM products WHERE status = 'ACTIVE'
CREATE INDEX idx_products_status ON products(status);

-- 5. Привязываем нашу функцию-триггер к таблице products
CREATE TRIGGER update_products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();