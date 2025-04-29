import React from 'react';
import EmployeeItem from './EmployeeItem';

const EmployeeList = ({ employees, onUpdate, onDelete, onEmployeeClick }) => {
  return (
    <section className="employee-list">
      {employees.length === 0 ? (
        <span className="no-employees-found">No employees found. Add a new employee to get started.</span>
      ) : (
        employees.map(employee => (
          <EmployeeItem
            key={employee.id}
            employee={employee}
            onUpdate={onUpdate}
            onDelete={onDelete}
            onEmployeeClick={onEmployeeClick}
          />
        ))
      )}
    </section>
  );
};

export default EmployeeList;
