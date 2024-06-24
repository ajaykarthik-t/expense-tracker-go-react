import React, { useState, useEffect } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faTrash } from '@fortawesome/free-solid-svg-icons';
import {
  PieChart,
  Pie,
  Tooltip,
  Legend,
  Cell,
} from 'recharts';

const Transactions = () => {
  const [title, setTitle] = useState('');
  const [amount, setAmount] = useState('');
  const [description, setDescription] = useState('');
  const [type, setType] = useState('income');
  const [message, setMessage] = useState('');
  const [transactions, setTransactions] = useState([]);
  const [totalIncome, setTotalIncome] = useState(0);
  const [totalExpense, setTotalExpense] = useState(0);

  useEffect(() => {
    fetchTransactions();
  }, []);

  useEffect(() => {
    calculateTotals();
  }, [transactions]);

  const calculateTotals = () => {
    let incomeTotal = 0;
    let expenseTotal = 0;

    transactions.forEach(trans => {
      if (trans.type === 'income') {
        incomeTotal += parseFloat(trans.amount);
      } else if (trans.type === 'expense') {
        expenseTotal += parseFloat(trans.amount);
      }
    });

    setTotalIncome(incomeTotal);
    setTotalExpense(expenseTotal);
  };

  const createTransaction = () => {
    const token = localStorage.getItem('token');
    if (!token) {
      setMessage('User is not authenticated. Please log in first.');
      return;
    }

    fetch('http://127.0.0.1:5000/transactions', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + token
      },
      body: JSON.stringify({ title, amount: parseFloat(amount), description, type })
    })
    .then(response => response.json())
    .then(data => {
      if (data.error) {
        setMessage(data.error);
      } else {
        setMessage('Transaction created');
        fetchTransactions();
      }
    })
    .catch(error => {
      setMessage('An error occurred: ' + error.message);
    });
  };

  const fetchTransactions = () => {
    const token = localStorage.getItem('token');
    fetch('http://127.0.0.1:5000/transactions', {
      method: 'GET',
      headers: {
        'Authorization': 'Bearer ' + token
      }
    })
    .then(response => response.json())
    .then(data => {
      if (data.error) {
        setMessage(data.error);
      } else {
        setTransactions(data);
      }
    })
    .catch(error => {
      setMessage('An error occurred: ' + error.message);
    });
  };

  const deleteTransaction = (id) => {
    const token = localStorage.getItem('token');
    if (!token) {
      setMessage('User is not authenticated. Please log in first.');
      return;
    }
  
    fetch(`http://127.0.0.1:5000/transactions/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': 'Bearer ' + token
      }
    })
    .then(response => {
      if (response.ok) {
        setMessage('Transaction deleted');
        // Filter out the deleted transaction from state
        setTransactions(transactions.filter(trans => trans.id !== id));
      } else {
        setMessage('Failed to delete transaction');
      }
    })
    .catch(error => {
      setMessage('An error occurred: ' + error.message);
    });
  };

  const pieChartData = [
    { name: 'Income', value: totalIncome },
    { name: 'Expense', value: totalExpense },
  ];

  const COLORS = ['#0088FE', '#FF8042'];

  return (
    <div className="container">
      <h2>Create Transaction</h2>
      <input
        type="text"
        placeholder="Title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
      />
      <input
        type="number"
        placeholder="Amount"
        value={amount}
        onChange={(e) => setAmount(e.target.value)}
      />
      <input
        type="text"
        placeholder="Description"
        value={description}
        onChange={(e) => setDescription(e.target.value)}
      />
      <select value={type} onChange={(e) => setType(e.target.value)}>
        <option value="income">Income</option>
        <option value="expense">Expense</option>
      </select>
      <button onClick={createTransaction}>Create</button>
      <div className="message">{message}</div>
  
      <h2>Transactions</h2>
      <button onClick={fetchTransactions}>Get Transactions</button>
      <ul className="transactions-list">
        {transactions.map((trans, index) => (
          <li key={index} className="transactions-item">
            <div className="transaction-info">
              <div className="transactions-item-title">{trans.title}</div>
              <div className="transactions-item-amount">${trans.amount}</div>
              <div className={`transactions-item-type ${trans.type}`}>
                {trans.type}
              </div>
              <div className="transactions-item-description">
                {trans.description}
              </div>
              <FontAwesomeIcon icon={faTrash} onClick={() => deleteTransaction(trans.id)} className="delete-icon" />
            </div>
          </li>
        ))}
      </ul>
  
      <div className="summary-container">
        <h1 className="summary-title">Summary</h1>
        <div className="summary-item">
          <label>Total Income:</label>
          <div className={`amount ${totalIncome >= 0 ? '' : 'negative'}`}>
            ${totalIncome}
          </div>
        </div>
        <div className="summary-item">
          <label>Total Expense:</label>
          <div className={`amount ${totalExpense >= 0 ? '' : 'negative'}`}>
            ${totalExpense}
          </div>
        </div>
        <div className="summary-item">
          <label>Account Balance:</label>
          <div className={`amount ${totalIncome - totalExpense >= 0 ? '' : 'negative'}`}>
            ${totalIncome - totalExpense}
          </div>
        </div>
      </div>

      <div className="pie-chart-container">
        <h2>Income vs Expense</h2>
        <PieChart width={400} height={300}>
          <Pie
            data={pieChartData}
            dataKey="value"
            nameKey="name"
            cx="50%"
            cy="50%"
            outerRadius={80}
            fill="#8884d8"
            label
          >
            {
              pieChartData.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
              ))
            }
          </Pie>
          <Tooltip />
          <Legend />
        </PieChart>
      </div>
    </div>
  );
};

export default Transactions;
