// TODO: Implement Transactions page
//
// This page should:
// 1. Display a form to create new transactions (select device, input data)
// 2. Show the resulting signature after signing
// 3. Display all transactions in a table with filtering by device
// 4. Show transaction details: counter, timestamp, data preview, signature
// 5. Handle loading and error states
//
// Suggested components:
// - TransactionForm component for signing data
// - TransactionTable component to display transactions

export default function TransactionsPage() {
  return (
    <div className="min-h-screen p-8">
      <div className="max-w-6xl mx-auto">
        <h1 className="text-3xl font-bold mb-8">Transactions</h1>
        
        <div className="mb-8">
          <p className="text-gray-600">
            TODO: Implement transaction form and list
          </p>
        </div>

        {/* TODO: Add TransactionForm component */}
        {/* TODO: Add TransactionTable component */}
      </div>
    </div>
  );
}
