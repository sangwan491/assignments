import React, { useState, useEffect } from 'react';

const EmployeeForm = ({ employees, onSave, onCancel }) => {
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    phone: '',
    department: '',
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  return (
    <form className="employee-form">
      <section className="form-group">
        <label>Name:</label>
        <input
          type="text"
          name="name"
          value={formData.name}
          onChange={handleChange}
          required
        />
      </section>
      
      <section className="form-group">
        <label>Email:</label>
        <input
          type="email"
          name="email"
          value={formData.email}
          onChange={handleChange}
          required
        />
      </section>
      
      <section className="form-group">
        <label>Phone:</label>
        <input
          type="tel"
          name="phone"
          value={formData.phone}
          onChange={handleChange}
          required
        />
      </section>
      
      <section className="form-group">
        <label>Department:</label>
        <input
          type="text"
          name="department"
          value={formData.department}
          onChange={handleChange}
          required
        />
      </section>
      
      <section className="form-actions">
        <button type="button" onClick={() => onSave(formData)}>Save</button>
        <button type="button" onClick={onCancel}>Cancel</button>
      </section>
    </form>
  );
};

export default EmployeeForm;