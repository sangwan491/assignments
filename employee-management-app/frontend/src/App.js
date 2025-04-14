import React, { useState, useEffect } from 'react';
import SearchBar from './components/SearchBar';
import EmployeeForm from './components/EmployeeForm';
import EmployeeList from './components/EmployeeList';
import './App.css';

import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080/api';

function App() {
  const [value, setValue] = useState("");
  const [isAdding, setAdding] = useState(false);
  const [employees, setEmployees] = useState([]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    fetchEmployees();
  }, []);

  const fetchEmployees = async () => {
    setIsLoading(true);
    
    try {
      const response = await fetch(`${API_BASE_URL}/employees`);
      
      const jsonResponse = await response.json();

      if (!response.ok) {
        throw new Error(jsonResponse.error);
      }
      
      setEmployees(jsonResponse);
    } catch (err) {
      console.log(err);
      toast.error(`Failed to load employees: ${err.message}`);
    } finally {
      setIsLoading(false);
    }
  };


  const onSave = async (newEmployee) => {
    setIsLoading(true);
    
    try {
      const response = await fetch(`${API_BASE_URL}/employees`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newEmployee),
      });
      
      const jsonResponse = await response.json();
      
      if (!response.ok) {
        throw new Error(jsonResponse.error);
      }

      newEmployee = jsonResponse.data;

      setEmployees((prevEmployees) => [...prevEmployees, newEmployee]);
      
      toast.success(jsonResponse.message);
      setAdding(false);  
    } catch (err) {
      console.log(err);
      toast.error(`Failed to add employee: ${err.message}`);
    } finally {
      setIsLoading(false);
    }
  };
  
  const onCancel = () => {
    setAdding(false);
  };


  const onUpdate = async (updatedEmployee) => {
    setIsLoading(true);
    
    try {
      const response = await fetch(`${API_BASE_URL}/employees/${updatedEmployee.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(updatedEmployee),
      });
      
      const jsonResponse = await response.json();

      if (!response.ok) {
        throw new Error(jsonResponse.error);
      }

      setEmployees((prevEmployees) => {
        return prevEmployees.map(employee => 
          employee.id === updatedEmployee.id ? updatedEmployee: employee
        );
      });

      toast.success(jsonResponse.message);
    } catch (err) {
      console.log(`err: ${err}`);
      toast.error(`Failed to update employee: ${err.message}`);
    } finally {
      setIsLoading(false);
    }
  };


  const onDelete = async (id) => {
    toast.promise(
      deleteEmployee(id),
      {
        pending: 'Deleting employee...',
        success: 'Employee deleted successfully',
        error: {
          render({data}) {
            return data.message;
          }
        }
      }
    );
  };
  
  const deleteEmployee = async (id) => {
    setIsLoading(true);
    
    try {
      const response = await fetch(`${API_BASE_URL}/employees/${id}`, {
        method: 'DELETE',
      });
      
      const jsonResponse = await response.json();

      if (!response.ok) {
        throw new Error(jsonResponse.error);
      }
      
      setEmployees((prevEmployees) =>
        prevEmployees.filter(employee => employee.id !== id)
      );

      return jsonResponse;
    } catch (err) {
      console.log(err);
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  const filteredEmployees = employees.filter(employee => {
    if (!value) return true;
    return JSON.stringify(employee).toLowerCase().includes(value.toLowerCase());
  });

  return (
    <>      
      <header>
        Employee Management
      </header>
      
      <main className="app-content">
      <ToastContainer
        position="top-right"
        autoClose={3000}
        hideProgressBar={false}
        newestOnTop
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
      />

        {isAdding ? (
          <EmployeeForm 
            employees={employees} 
            setEmployees={setEmployees} 
            onSave={onSave} 
            onCancel={onCancel}
            isLoading={isLoading}
          />
        ) : (
          <button 
            className='add-new-btn' 
            onClick={() => setAdding(true)}
            disabled={isLoading}
          >
            Add new Employee
          </button>
        )}
        
        <SearchBar value={value} setValue={setValue} />
        
        {isLoading? (
          <span className="loading">Loading employees...</span>
        ) : (
          <EmployeeList
            employees={filteredEmployees}
            onUpdate={onUpdate}
            onDelete={onDelete}
            isLoading={isLoading}
          />
        )}
      </main>
    </>
  );
}

export default App;