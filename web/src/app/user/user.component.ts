import { Component, OnInit } from '@angular/core';
import { BackendService, UserItem } from '../backend.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.scss'],
})
export class UserComponent implements OnInit {
  private id: string = '';
  public user: UserDisplay | undefined;

  constructor(
    private backendService: BackendService,
    private router: ActivatedRoute,
  ) {}

  ngOnInit(): void {
    this.id = String(this.router.snapshot.paramMap.get('id'));

    this.backendService.getUser(this.id).subscribe((response) => {
      this.user = {
        id: response.data.id,
        username: response.data.attributes.username,
        created_at: new Date(response.data.attributes.created_at),
        updated_at: new Date(response.data.attributes.updated_at),
        active: response.data.attributes.active,
      };
    });
  }
}

export interface UserDisplay {
  id: string;
  username: string;
  created_at: Date;
  updated_at: Date;
  active: boolean;
}
