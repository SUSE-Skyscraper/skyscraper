import { Component, Input, OnChanges, OnInit } from '@angular/core';
import { BackendService, UserItem } from '../backend.service';

@Component({
  selector: 'app-audit-log',
  templateUrl: './audit-log.component.html',
  styleUrls: ['./audit-log.component.scss'],
})
export class AuditLogComponent implements OnInit, OnChanges {
  auditLogs: AuditLogTableItem[] = [];
  auditLogColumns: string[] = [];

  @Input() resource_id?: string;
  @Input() resource_type?: string;

  constructor(private backendService: BackendService) {}

  ngOnInit(): void {
    this.getAuditLogs();
  }

  ngOnChanges() {
    this.getAuditLogs();
  }

  private getAuditLogs() {
    if (this.resource_id !== undefined && this.resource_type !== undefined) {
      this.auditLogColumns = ['user_id', 'message', 'created_at'];
    } else {
      this.auditLogColumns = [
        'user_id',
        'resource_type',
        'resource_id',
        'message',
        'created_at',
      ];
    }
    this.backendService
      .getAuditLogs(this.resource_id, this.resource_type)
      .subscribe((response) => {
        if (response.data !== null) {
          let users = new Map<string, UserItem>();

          response.included?.forEach((item) => {
            if (item.type === 'user') {
              users.set(item.id, item as UserItem);
            }
          });

          let tableItems: AuditLogTableItem[] = [];
          response.data.forEach((item) => {
            const user = users.get(item.relationships['user']?.data.id);

            tableItems.push({
              user_id: item.attributes.user_id,
              message: item.attributes.message,
              username: user?.attributes.username ?? '',
              created_at: new Date(item.attributes.created_at),
              resource_id: item.attributes.resource_id,
              resource_type: item.attributes.resource_type,
            });
          });

          this.auditLogs = tableItems;
        } else {
          this.auditLogs = [];
        }
      });
  }
}

export interface AuditLogTableItem {
  message: string;
  user_id: string;
  username: string;
  created_at: Date;
  resource_type: string;
  resource_id: string;
}
