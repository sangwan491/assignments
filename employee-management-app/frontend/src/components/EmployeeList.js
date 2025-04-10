import React from 'react';
import EmployeeItem from './EmployeeItem';

const EmployeeList = ({ employees, onUpdate, onDelete }) => {
  console.log(employees);

  return (
    <div className="employee-list">
      {employees.length === 0 ? (
        <p>No employees found. Add a new employee to get started.</p>
      ) : (
        employees.map(employee => (
          <EmployeeItem
            key={employee.id}
            employee={employee}
            onUpdate={onUpdate}
            onDelete={onDelete}
          />
        ))
      )}
    </div>
  );
};

export default EmployeeList;
