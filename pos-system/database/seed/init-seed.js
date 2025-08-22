// Seed data for POS System
// This file creates initial data for development and testing

db = db.getSiblingDB('pos_db');

// Create admin user collection with initial admin user
db.users.insertMany([
  {
    email: 'admin@posstore.com',
    name: 'System Administrator',
    role: 'ADMIN',
    isActive: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    email: 'manager@posstore.com',
    name: 'Store Manager',
    role: 'MANAGER',
    isActive: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    email: 'cashier@posstore.com',
    name: 'Store Cashier',
    role: 'CASHIER',
    isActive: true,
    createdAt: new Date(),
    updatedAt: new Date()
  }
]);

// Create initial product categories
db.categories.insertMany([
  {
    name: 'Electronics',
    description: 'Electronic devices and accessories',
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: 'Clothing',
    description: 'Apparel and fashion items',
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: 'Food & Beverages',
    description: 'Food items and drinks',
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: 'Books',
    description: 'Books and publications',
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: 'Home & Garden',
    description: 'Home improvement and garden supplies',
    createdAt: new Date(),
    updatedAt: new Date()
  }
]);

// Get category IDs for product insertion
const electronics = db.categories.findOne({name: 'Electronics'});
const clothing = db.categories.findOne({name: 'Clothing'});
const food = db.categories.findOne({name: 'Food & Beverages'});
const books = db.categories.findOne({name: 'Books'});
const home = db.categories.findOne({name: 'Home & Garden'});

// Create initial products
db.products.insertMany([
  // Electronics
  {
    name: 'Wireless Bluetooth Headphones',
    description: 'High-quality wireless headphones with noise cancellation',
    sku: 'ELC-WBH-001',
    barcode: '1234567890123',
    price: 79.99,
    cost: 45.00,
    stock: 25,
    minStock: 5,
    maxStock: 50,
    categoryId: electronics._id,
    isActive: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: 'USB-C Charging Cable',
    description: 'Fast charging USB-C cable 6ft length',
    sku: 'ELC-USB-002',
    barcode: '1234567890124',
    price: 19.99,
    cost: 8.00,
    stock: 50,
    minStock: 10,
    maxStock: 100,
    categoryId: electronics._id,
    isActive: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  
  // Clothing
  {
    name: 'Cotton T-Shirt',
    description: 'Comfortable 100% cotton t-shirt',
    sku: 'CLO-TSH-001',
    barcode: '1234567890125',
    price: 24.99,
    cost: 12.00,
    stock: 30,
    minStock: 5,
    maxStock: 60,
    categoryId: clothing._id,
    isActive: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: 'Denim Jeans',
    description: 'Classic fit denim jeans',
    sku: 'CLO-JEA-002',
    barcode: '1234567890126',
    price: 59.99,
    cost: 30.00,
    stock: 20,
    minStock: 3,
    maxStock: 40,
    categoryId: clothing._id,
    isActive: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  
  // Food & Beverages
  {
    name: 'Organic Coffee Beans',
    description: 'Premium organic coffee beans 1lb bag',
    sku: 'FOD-COF-001',
    barcode: '1234567890127',
    price: 15.99,
    cost: 8.50,
    stock: 40,
    minStock: 8,
    maxStock: 80,
    categoryId: food._id,
    isActive: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: 'Sparkling Water 12-pack',
    description: 'Natural sparkling water 12 cans',
    sku: 'FOD-WAT-002',
    barcode: '1234567890128',
    price: 8.99,
    cost: 5.00,
    stock: 35,
    minStock: 6,
    maxStock: 70,
    categoryId: food._id,
    isActive: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  
  // Books
  {
    name: 'JavaScript Programming Guide',
    description: 'Complete guide to modern JavaScript programming',
    sku: 'BOO-JAV-001',
    barcode: '1234567890129',
    price: 39.99,
    cost: 20.00,
    stock: 15,
    minStock: 3,
    maxStock: 30,
    categoryId: books._id,
    isActive: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  
  // Home & Garden
  {
    name: 'Indoor Plant Pot',
    description: 'Ceramic plant pot for indoor plants',
    sku: 'HOM-POT-001',
    barcode: '1234567890130',
    price: 12.99,
    cost: 6.00,
    stock: 22,
    minStock: 4,
    maxStock: 50,
    categoryId: home._id,
    isActive: true,
    createdAt: new Date(),
    updatedAt: new Date()
  }
]);

// Create initial expense categories for reference
db.expenses.insertMany([
  {
    description: 'Monthly Store Rent',
    amount: 2500.00,
    category: 'rent',
    vendor: 'Property Management Co.',
    date: new Date(),
    isRecurring: true,
    userId: db.users.findOne({email: 'admin@posstore.com'})._id,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    description: 'Electricity Bill',
    amount: 150.00,
    category: 'utilities',
    vendor: 'Power Company',
    date: new Date(),
    isRecurring: true,
    userId: db.users.findOne({email: 'admin@posstore.com'})._id,
    createdAt: new Date(),
    updatedAt: new Date()
  }
]);

print('Seed data inserted successfully!');
print('Created users:', db.users.count());
print('Created categories:', db.categories.count());
print('Created products:', db.products.count());
print('Created expenses:', db.expenses.count());
