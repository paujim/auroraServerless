import React from 'react';
import MaterialTable from 'material-table';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import axios from 'axios'

export default function ProfilesTable(prop) {

  const showError = prop.showError;

  return (
    <MaterialTable
      title="Profiles"
      options={{
        search: false
      }}
      columns={[
        {
          title: 'Full Name',
          field: 'full_name'
        },
        {
          title: 'Email',
          field: 'email'
        },
        {
          title: 'Phone Numbers',
          field: 'phone_numbers',
          render: rowData => {
            return (
              <List>{rowData.phone_numbers.map(phone => (
                <ListItem>
                  <ListItemText primary={phone} />
                </ListItem>
              ))}
              </List>)
          },
        },
      ]}
      data={query =>
        new Promise((resolve, reject) => {
          let url = 'http://localhost:5000/'
          axios.get(url)
            .then(response => {
              let data = response.data
              if (response.status !== 200) {
                throw Error(data.error);
              }
              resolve({
                data: data.profiles,
                page: 0,
                totalCount: 1,
              })
            })
            .catch(error => {
              showError(error.message)
              resolve({
                data: [],
                page: 0,
                totalCount: 1,
              })
            })
        })
      }
      editable={{
        isEditable: (row) => false,
        onRowAdd: (newData) =>
          new Promise((resolve, reject) => {
            let phoneNumbers = newData.phone_numbers.split(";")
            newData.phone_numbers = phoneNumbers;
            let url = 'http://localhost:5000/'
            axios.post(url, JSON.stringify(newData))
              .then(response => {
                resolve()
              })
              .catch(error => {
                if (error.response) {
                  showError(error.response.data.error)
                }
                resolve()
              })
          }),
        onRowUpdate: (newData, oldData) => alert("the update endpoint is not implemented yet"),
      }}
    />
  );
}