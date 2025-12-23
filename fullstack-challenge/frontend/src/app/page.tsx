import Link from "next/link";

export default function Home() {
  return (
    <div className="min-h-screen p-8 pb-20 gap-16 sm:p-20">
      <main className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold mb-8">Signature Service</h1>
        
        <div className="grid gap-4 mb-8">
          <Link 
            href="/devices"
            className="p-6 border border-gray-300 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
          >
            <h2 className="text-2xl font-semibold mb-2">Devices →</h2>
            <p className="text-gray-600 dark:text-gray-400">
              Manage signature devices and create new ones
            </p>
          </Link>

          <Link 
            href="/transactions"
            className="p-6 border border-gray-300 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
          >
            <h2 className="text-2xl font-semibold mb-2">Transactions →</h2>
            <p className="text-gray-600 dark:text-gray-400">
              Sign data and view transaction history
            </p>
          </Link>
        </div>

        <div className="mt-16 text-sm text-gray-500">
          <p>fiskaly Fullstack Coding Challenge</p>
        </div>
      </main>
    </div>
  );
}
