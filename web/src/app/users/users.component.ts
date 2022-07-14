import { Component, OnInit } from '@angular/core';
import { BackendService, UserItem } from '../backend.service';

@Component({
  selector: 'app-users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.scss'],
})
export class UsersComponent implements OnInit {
  public users: UserTableItem[] = [];
  public userTableColumns: string[] = [
    'username',
    'active',
    'created_at',
    'updated_at',
  ];

  constructor(private backendService: BackendService) {}

  ngOnInit(): void {
    this.backendService.getUsers().subscribe((response) => {
      let users: UserTableItem[] = [];
      response.data.forEach((item) => {
        users.push({
          id: item.id,
          username: item.attributes.username,
          created_at: new Date(item.attributes.created_at),
          updated_at: new Date(item.attributes.updated_at),
          active: item.attributes.active,
        });
      });
      this.users = users;
    });
  }
}

export interface UserTableItem {
  id: string;
  username: string;
  created_at: Date;
  updated_at: Date;
  active: boolean;
}
