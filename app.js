import React, { useState, useEffect, createContext, useContext } from 'react';
import { initializeApp } from 'firebase/app';
import { getAuth, signInAnonymously, signInWithCustomToken, onAuthStateChanged } from 'firebase/auth';
import { getFirestore, collection, query, orderBy, onSnapshot, doc, addDoc, setDoc, deleteDoc, getDocs } from 'firebase/firestore';

// Context for Firebase and User Auth
const AppContext = createContext(null);

// Utility function for generating unique IDs (if needed for anonymous users)
const generateUniqueId = () => crypto.randomUUID();

// Main App Component
function App() {
  const [db, setDb] = useState(null);
  const [auth, setAuth] = useState(null);
  const [userId, setUserId] = useState(null);
  const [isAuthReady, setIsAuthReady] = useState(false);
  const [currentPage, setCurrentPage] = useState('trades'); // 'trades', 'cost-basis', 'fif-calculation', 'fif-report'

  useEffect(() => {
    // Initialize Firebase
    const appId = typeof __app_id !== 'undefined' ? __app_id : 'default-app-id';
    const firebaseConfig = typeof __firebase_config !== 'undefined' ? JSON.parse(__firebase_config) : {};

    try {
      const app = initializeApp(firebaseConfig);
      const firestoreDb = getFirestore(app);
      const firebaseAuth = getAuth(app);

      setDb(firestoreDb);
      setAuth(firebaseAuth);

      // Listen for authentication state changes
      const unsubscribe = onAuthStateChanged(firebaseAuth, async (user) => {
        if (user) {
          setUserId(user.uid);
        } else {
          // Sign in anonymously if no initial token is provided
          if (typeof __initial_auth_token !== 'undefined' && __initial_auth_token) {
            try {
              await signInWithCustomToken(firebaseAuth, __initial_auth_token);
              setUserId(firebaseAuth.currentUser?.uid); // Ensure userId is set after successful sign-in
            } catch (error) {
              console.error("Error signing in with custom token:", error);
              await signInAnonymously(firebaseAuth);
              setUserId(firebaseAuth.currentUser?.uid || generateUniqueId());
            }
          } else {
            await signInAnonymously(firebaseAuth);
            setUserId(firebaseAuth.currentUser?.uid || generateUniqueId());
          }
        }
        setIsAuthReady(true); // Mark auth as ready once initial check is done
      });

      return () => unsubscribe(); // Cleanup auth listener on unmount
    } catch (error) {
      console.error("Error initializing Firebase:", error);
      // Handle Firebase initialization error, e.g., display a message to the user
    }
  }, []); // Empty dependency array ensures this runs only once on mount

  if (!isAuthReady) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-100">
        <div className="text-lg font-semibold text-gray-700">Loading application...</div>
      </div>
    );
  }

  return (
    <AppContext.Provider value={{ db, auth, userId, isAuthReady }}>
      <div className="min-h-screen bg-gray-50 font-sans antialiased text-gray-800 flex flex-col">
        {/* Header */}
        <Header userId={userId} />

        {/* Main Content Area */}
        <div className="flex flex-1 flex-col lg:flex-row">
          {/* Navigation Sidebar/Tabs */}
          <Navigation currentPage={currentPage} setCurrentPage={setCurrentPage} />

          {/* Page Content */}
          <main className="flex-1 p-4 md:p-8 lg:p-10 overflow-auto">
            {currentPage === 'trades' && <TradesPage />}
            {currentPage === 'cost-basis' && <CostBasisPage />}
            {currentPage === 'fif-calculation' && <FifCalculationPage />}
            {currentPage === 'fif-report' && <FifReportPage />}
          </main>
        </div>
      </div>
    </AppContext.Provider>
  );
}

// Header Component
const Header = ({ userId }) => {
  return (
    <header className="bg-gradient-to-r from-blue-600 to-indigo-700 text-white p-4 shadow-lg z-10">
      <div className="container mx-auto flex flex-col sm:flex-row justify-between items-center">
        <h1 className="text-3xl font-bold tracking-tight mb-2 sm:mb-0">FIF Tax Calculator</h1>
        {userId && (
          <div className="text-sm bg-blue-700 px-3 py-1 rounded-full opacity-90">
            User ID: <span className="font-mono text-blue-200">{userId}</span>
          </div>
        )}
      </div>
    </header>
  );
};

// Navigation Component
const Navigation = ({ currentPage, setCurrentPage }) => {
  const navItems = [
    { id: 'trades', name: 'Manage Trades' },
    { id: 'cost-basis', name: 'Cost Basis' },
    { id: 'fif-calculation', name: 'FIF Calculation' },
    { id: 'fif-report', name: 'FIF Report' },
  ];

  return (
    <nav className="bg-white shadow-md lg:w-64 p-4 lg:p-6 border-b lg:border-r border-gray-200 flex flex-row lg:flex-col overflow-x-auto lg:overflow-x-hidden">
      <ul className="flex flex-row lg:flex-col space-x-4 lg:space-x-0 lg:space-y-2 w-full">
        {navItems.map((item) => (
          <li key={item.id} className="flex-shrink-0">
            <button
              onClick={() => setCurrentPage(item.id)}
              className={`block w-full text-left py-2 px-4 rounded-lg transition-all duration-200 ease-in-out
                ${currentPage === item.id
                  ? 'bg-blue-500 text-white shadow-md'
                  : 'text-gray-700 hover:bg-gray-100 hover:text-blue-600'
                }`}
            >
              {item.name}
            </button>
          </li>
        ))}
      </ul>
    </nav>
  );
};

// Trades Page Component
const TradesPage = () => {
  const { db, userId, isAuthReady } = useContext(AppContext);
  const [trades, setTrades] = useState([]);
  const [newTrade, setNewTrade] = useState({
    ticker: '',
    date: '',
    type: 'Buy',
    quantity: '',
    price: '',
    currency: 'USD',
  });
  const [editingTradeId, setEditingTradeId] = useState(null);
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState(''); // 'success' or 'error'

  const tradesCollectionRef = db && userId ? collection(db, `artifacts/${__app_id}/users/${userId}/trades`) : null;

  useEffect(() => {
    if (!tradesCollectionRef || !isAuthReady) return;

    const q = query(tradesCollectionRef); // No orderBy to avoid index issues

    const unsubscribe = onSnapshot(q, (snapshot) => {
      const fetchedTrades = snapshot.docs.map(doc => ({
        id: doc.id,
        ...doc.data()
      }));
      setTrades(fetchedTrades);
    }, (error) => {
      console.error("Error fetching trades:", error);
      showMessage("Error fetching trades. Please try again.", "error");
    });

    return () => unsubscribe();
  }, [tradesCollectionRef, isAuthReady]);

  const showMessage = (text, type) => {
    setMessage(text);
    setMessageType(type);
    setTimeout(() => {
      setMessage('');
      setMessageType('');
    }, 3000); // Message disappears after 3 seconds
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setNewTrade(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!tradesCollectionRef) {
      showMessage("Database not ready. Please try again.", "error");
      return;
    }

    // Basic validation
    if (!newTrade.ticker || !newTrade.date || !newTrade.quantity || !newTrade.price) {
      showMessage("Please fill in all required fields.", "error");
      return;
    }
    if (isNaN(newTrade.quantity) || isNaN(newTrade.price) || parseFloat(newTrade.quantity) <= 0 || parseFloat(newTrade.price) <= 0) {
      showMessage("Quantity and Price must be positive numbers.", "error");
      return;
    }

    try {
      const tradeData = {
        ...newTrade,
        quantity: parseFloat(newTrade.quantity),
        price: parseFloat(newTrade.price),
        timestamp: new Date().toISOString(), // Store creation timestamp
      };

      if (editingTradeId) {
        await setDoc(doc(tradesCollectionRef, editingTradeId), tradeData);
        showMessage("Trade updated successfully!", "success");
        setEditingTradeId(null);
      } else {
        await addDoc(tradesCollectionRef, tradeData);
        showMessage("Trade added successfully!", "success");
      }
      setNewTrade({
        ticker: '',
        date: '',
        type: 'Buy',
        quantity: '',
        price: '',
        currency: 'USD',
      });
    } catch (error) {
      console.error("Error adding/updating trade:", error);
      showMessage("Error saving trade. Please try again.", "error");
    }
  };

  const handleEdit = (trade) => {
    setEditingTradeId(trade.id);
    setNewTrade({
      ticker: trade.ticker,
      date: trade.date,
      type: trade.type,
      quantity: trade.quantity,
      price: trade.price,
      currency: trade.currency,
    });
    showMessage("Editing trade...", "success");
  };

  const handleDelete = async (id) => {
    if (!tradesCollectionRef) {
      showMessage("Database not ready. Please try again.", "error");
      return;
    }
    if (window.confirm("Are you sure you want to delete this trade?")) { // Using window.confirm for simplicity, replace with custom modal if needed
      try {
        await deleteDoc(doc(tradesCollectionRef, id));
        showMessage("Trade deleted successfully!", "success");
      } catch (error) {
        console.error("Error deleting trade:", error);
        showMessage("Error deleting trade. Please try again.", "error");
      }
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-lg">
      <h2 className="text-2xl font-semibold mb-6 text-gray-900">Manage Your Trades</h2>

      {message && (
        <div className={`p-3 mb-4 rounded-md text-sm ${messageType === 'success' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
          {message}
        </div>
      )}

      {/* Trade Input Form */}
      <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-8 p-4 border border-gray-200 rounded-md bg-gray-50">
        <div className="flex flex-col">
          <label htmlFor="ticker" className="block text-sm font-medium text-gray-700 mb-1">Stock Ticker</label>
          <input
            type="text"
            id="ticker"
            name="ticker"
            value={newTrade.ticker}
            onChange={handleChange}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 p-2"
            placeholder="e.g., AAPL"
            required
          />
        </div>
        <div className="flex flex-col">
          <label htmlFor="date" className="block text-sm font-medium text-gray-700 mb-1">Date</label>
          <input
            type="date"
            id="date"
            name="date"
            value={newTrade.date}
            onChange={handleChange}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 p-2"
            required
          />
        </div>
        <div className="flex flex-col">
          <label htmlFor="type" className="block text-sm font-medium text-gray-700 mb-1">Type</label>
          <select
            id="type"
            name="type"
            value={newTrade.type}
            onChange={handleChange}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 p-2"
          >
            <option value="Buy">Buy</option>
            <option value="Sell">Sell</option>
          </select>
        </div>
        <div className="flex flex-col">
          <label htmlFor="quantity" className="block text-sm font-medium text-gray-700 mb-1">Quantity</label>
          <input
            type="number"
            id="quantity"
            name="quantity"
            value={newTrade.quantity}
            onChange={handleChange}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 p-2"
            placeholder="e.g., 100"
            min="0.01"
            step="any"
            required
          />
        </div>
        <div className="flex flex-col">
          <label htmlFor="price" className="block text-sm font-medium text-gray-700 mb-1">Price per Share</label>
          <input
            type="number"
            id="price"
            name="price"
            value={newTrade.price}
            onChange={handleChange}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 p-2"
            placeholder="e.g., 150.25"
            min="0.01"
            step="any"
            required
          />
        </div>
        <div className="flex flex-col">
          <label htmlFor="currency" className="block text-sm font-medium text-gray-700 mb-1">Currency</label>
          <select
            id="currency"
            name="currency"
            value={newTrade.currency}
            onChange={handleChange}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 p-2"
          >
            <option value="USD">USD</option>
            <option value="NZD">NZD</option>
            <option value="AUD">AUD</option>
            {/* Add more currencies as needed */}
          </select>
        </div>
        <div className="md:col-span-2 lg:col-span-3 flex justify-end">
          <button
            type="submit"
            className="inline-flex justify-center py-2 px-6 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors duration-200"
          >
            {editingTradeId ? 'Update Trade' : 'Add Trade'}
          </button>
        </div>
      </form>

      {/* Trades Table */}
      <h3 className="text-xl font-semibold mb-4 text-gray-900">Your Trades</h3>
      {trades.length === 0 ? (
        <p className="text-gray-600">No trades added yet. Use the form above to add your first trade.</p>
      ) : (
        <div className="overflow-x-auto rounded-lg shadow-md border border-gray-200">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-100">
              <tr>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider rounded-tl-lg">Ticker</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Quantity</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Price</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Currency</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider rounded-tr-lg">Actions</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {trades.map((trade) => (
                <tr key={trade.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{trade.ticker}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{trade.date}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${trade.type === 'Buy' ? 'bg-blue-100 text-blue-800' : 'bg-red-100 text-red-800'}`}>
                      {trade.type}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{trade.quantity}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{trade.price.toFixed(2)}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{trade.currency}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <button
                      onClick={() => handleEdit(trade)}
                      className="text-indigo-600 hover:text-indigo-900 mr-4 transition-colors duration-200"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => handleDelete(trade.id)}
                      className="text-red-600 hover:text-red-900 transition-colors duration-200"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

// Cost Basis Page Component
const CostBasisPage = () => {
  const { db, userId, isAuthReady } = useContext(AppContext);
  const [costBasisData, setCostBasisData] = useState({});
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState('');

  const tradesCollectionRef = db && userId ? collection(db, `artifacts/${__app_id}/users/${userId}/trades`) : null;

  useEffect(() => {
    if (!tradesCollectionRef || !isAuthReady) return;

    const q = query(tradesCollectionRef); // No orderBy to avoid index issues

    const unsubscribe = onSnapshot(q, (snapshot) => {
      const fetchedTrades = snapshot.docs.map(doc => ({
        id: doc.id,
        ...doc.data()
      }));
      calculateCostBasis(fetchedTrades);
    }, (error) => {
      console.error("Error fetching trades for cost basis:", error);
      showMessage("Error loading cost basis data.", "error");
    });

    return () => unsubscribe();
  }, [tradesCollectionRef, isAuthReady]);

  const showMessage = (text, type) => {
    setMessage(text);
    setMessageType(type);
    setTimeout(() => {
      setMessage('');
      setMessageType('');
    }, 3000);
  };

  const calculateCostBasis = (trades) => {
    const basis = {};

    trades.forEach(trade => {
      const { ticker, type, quantity, price } = trade;
      if (!basis[ticker]) {
        basis[ticker] = { totalQuantity: 0, totalCost: 0, currency: trade.currency };
      }

      if (type === 'Buy') {
        basis[ticker].totalQuantity += quantity;
        basis[ticker].totalCost += quantity * price;
      } else if (type === 'Sell') {
        // For simplicity, assuming FIFO for cost basis calculation here.
        // A more robust system would need to track individual lots.
        // For now, we'll just reduce quantity and adjust cost proportionally.
        // This is a simplified model and may not be accurate for all tax purposes.
        if (basis[ticker].totalQuantity > 0) {
          const costPerShare = basis[ticker].totalCost / basis[ticker].totalQuantity;
          basis[ticker].totalQuantity -= quantity;
          basis[ticker].totalCost -= quantity * costPerShare;
          // Ensure no negative quantities/costs
          if (basis[ticker].totalQuantity < 0) basis[ticker].totalQuantity = 0;
          if (basis[ticker].totalCost < 0) basis[ticker].totalCost = 0;
        }
      }
    });

    // Filter out stocks with zero quantity and cost (fully sold)
    const filteredBasis = Object.fromEntries(
      Object.entries(basis).filter(([, value]) => value.totalQuantity > 0)
    );

    setCostBasisData(filteredBasis);
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-lg">
      <h2 className="text-2xl font-semibold mb-6 text-gray-900">Cost Basis Overview</h2>

      {message && (
        <div className={`p-3 mb-4 rounded-md text-sm ${messageType === 'success' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
          {message}
        </div>
      )}

      {Object.keys(costBasisData).length === 0 ? (
        <p className="text-gray-600">No stocks currently held or no trades entered yet.</p>
      ) : (
        <div className="overflow-x-auto rounded-lg shadow-md border border-gray-200">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-100">
              <tr>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider rounded-tl-lg">Stock Ticker</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Total Quantity Held</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Total Cost</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Average Cost per Share</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider rounded-tr-lg">Currency</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {Object.entries(costBasisData).map(([ticker, data]) => (
                <tr key={ticker} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{ticker}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{data.totalQuantity.toFixed(2)}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{data.totalCost.toFixed(2)}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {(data.totalQuantity > 0 ? (data.totalCost / data.totalQuantity) : 0).toFixed(2)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{data.currency}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

// FIF Calculation Page Component
const FifCalculationPage = () => {
  const [selectedMethod, setSelectedMethod] = useState('FDR'); // Fair Dividend Rate (FDR) or Comparative Value (CV)
  const [financialYear, setFinancialYear] = useState(new Date().getFullYear()); // Current year as default
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState('');

  const showMessage = (text, type) => {
    setMessage(text);
    setMessageType(type);
    setTimeout(() => {
      setMessage('');
      setMessageType('');
    }, 3000);
  };

  const handleCalculate = () => {
    // This is where the actual FIF calculation logic would go.
    // It would involve fetching trades, applying the selected method (FDR/CV)
    // and then storing/displaying the result.
    // For this design, we'll just show a placeholder message.
    showMessage(`Calculating FIF tax for ${financialYear} using ${selectedMethod} method... (Logic to be implemented)`, "success");
    console.log(`Simulating calculation for year ${financialYear} with method ${selectedMethod}`);
  };

  const getYears = () => {
    const currentYear = new Date().getFullYear();
    const years = [];
    for (let i = currentYear; i >= currentYear - 10; i--) { // Last 10 years
      years.push(i);
    }
    return years;
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-lg">
      <h2 className="text-2xl font-semibold mb-6 text-gray-900">FIF Tax Calculation</h2>

      {message && (
        <div className={`p-3 mb-4 rounded-md text-sm ${messageType === 'success' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
          {message}
        </div>
      )}

      <div className="space-y-6">
        {/* Method Selection */}
        <div>
          <label className="block text-lg font-medium text-gray-700 mb-2">Select Calculation Method:</label>
          <div className="flex flex-col sm:flex-row space-y-3 sm:space-y-0 sm:space-x-6">
            <label className="inline-flex items-center">
              <input
                type="radio"
                name="fifMethod"
                value="FDR"
                checked={selectedMethod === 'FDR'}
                onChange={() => setSelectedMethod('FDR')}
                className="form-radio h-4 w-4 text-blue-600 transition-colors duration-200"
              />
              <span className="ml-2 text-gray-700 font-medium">Fair Dividend Rate (FDR)</span>
            </label>
            <label className="inline-flex items-center">
              <input
                type="radio"
                name="fifMethod"
                value="CV"
                checked={selectedMethod === 'CV'}
                onChange={() => setSelectedMethod('CV')}
                className="form-radio h-4 w-4 text-blue-600 transition-colors duration-200"
              />
              <span className="ml-2 text-gray-700 font-medium">Comparative Value (CV)</span>
            </label>
          </div>
          <p className="text-sm text-gray-500 mt-2">
            <strong>FDR:</strong> Assumes a 5% taxable income on the opening market value of your FIF investments. Often simpler.
          </p>
          <p className="text-sm text-gray-500 mt-1">
            <strong>CV:</strong> Calculates income based on the change in market value plus any distributions received, minus costs. More complex.
          </p>
        </div>

        {/* Financial Year Selection */}
        <div>
          <label htmlFor="financialYear" className="block text-lg font-medium text-gray-700 mb-2">Select Financial Year:</label>
          <select
            id="financialYear"
            name="financialYear"
            value={financialYear}
            onChange={(e) => setFinancialYear(parseInt(e.target.value))}
            className="mt-1 block w-full sm:w-1/2 lg:w-1/3 rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 p-2"
          >
            {getYears().map(year => (
              <option key={year} value={year}>{year}</option>
            ))}
          </select>
          <p className="text-sm text-gray-500 mt-2">
            New Zealand's financial year runs from 1 April to 31 March.
          </p>
        </div>

        {/* Calculate Button */}
        <div>
          <button
            onClick={handleCalculate}
            className="inline-flex justify-center py-3 px-8 border border-transparent shadow-sm text-base font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 transition-colors duration-200"
          >
            Calculate FIF Tax
          </button>
        </div>
      </div>
    </div>
  );
};

// FIF Report Page Component
const FifReportPage = () => {
  const [fifReports, setFifReports] = useState([]); // This would be populated by actual calculations
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState('');

  // Placeholder for fetching/displaying reports.
  // In a real app, this would fetch calculated reports from Firestore
  // or trigger the calculation if not already done.
  useEffect(() => {
    // Simulate fetching some dummy reports
    const dummyReports = [
      { year: 2023, method: 'FDR', fifIncome: 1200.50, fifTaxPayable: 396.17 },
      { year: 2022, method: 'CV', fifIncome: 850.00, fifTaxPayable: 280.50 },
      { year: 2021, method: 'FDR', fifIncome: 1500.00, fifTaxPayable: 495.00 },
    ];
    setFifReports(dummyReports);
  }, []);

  const showMessage = (text, type) => {
    setMessage(text);
    setMessageType(type);
    setTimeout(() => {
      setMessage('');
      setMessageType('');
    }, 3000);
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-lg">
      <h2 className="text-2xl font-semibold mb-6 text-gray-900">FIF Tax Reports</h2>

      {message && (
        <div className={`p-3 mb-4 rounded-md text-sm ${messageType === 'success' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
          {message}
        </div>
      )}

      {fifReports.length === 0 ? (
        <p className="text-gray-600">No FIF tax reports available yet. Please calculate your FIF tax on the "FIF Calculation" page.</p>
      ) : (
        <div className="overflow-x-auto rounded-lg shadow-md border border-gray-200">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-100">
              <tr>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider rounded-tl-lg">Financial Year</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Method Used</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">FIF Income (NZD)</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider rounded-tr-lg">FIF Tax Payable (NZD)</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {fifReports.map((report) => (
                <tr key={report.year} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{report.year}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{report.method}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{report.fifIncome.toFixed(2)}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{report.fifTaxPayable.toFixed(2)}</td>
                </tr>
              ))}
            </tbody>
          </table>
          <p className="p-4 text-sm text-gray-600 bg-gray-50 border-t border-gray-200 rounded-b-lg">
            *Please note: These are illustrative figures. Actual FIF tax calculations are complex and depend on various factors including foreign exchange rates, specific investment types, and IRD guidelines. Always consult with a qualified tax advisor.
          </p>
        </div>
      )}
    </div>
  );
};

export default App;
