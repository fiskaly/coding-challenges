import { migrate } from 'drizzle-orm/postgres-js/migrator';
import postgres from 'postgres';
import { drizzle } from 'drizzle-orm/postgres-js';

const runMigrations = async () => {
  const connectionString = process.env.DATABASE_URL || 
    'postgresql://postgres:postgres@localhost:5432/signature_service';
  
  const migrationClient = postgres(connectionString, { max: 1 });
  const db = drizzle(migrationClient);
  
  console.log('Running migrations...');
  await migrate(db, { migrationsFolder: './drizzle' });
  console.log('Migrations completed!');
  
  await migrationClient.end();
};

runMigrations().catch((err) => {
  console.error('Migration failed!', err);
  process.exit(1);
});
