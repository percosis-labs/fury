<!--
order: 4
-->

# Events

The jinx module emits the following events:

## Handlers

### MsgDeposit

| Type         | Attribute Key | Attribute Value       |
| ------------ | ------------- | --------------------- |
| message      | module        | jinx                  |
| message      | sender        | `{sender address}`    |
| jinx_deposit | amount        | `{amount}`            |
| jinx_deposit | depositor     | `{depositor address}` |

### MsgWithdraw

| Type            | Attribute Key | Attribute Value       |
| --------------- | ------------- | --------------------- |
| message         | module        | jinx                  |
| message         | sender        | `{sender address}`    |
| jinx_withdrawal | amount        | `{amount}`            |
| jinx_withdrawal | depositor     | `{depositor address}` |

### MsgBorrow

| Type            | Attribute Key | Attribute Value      |
| --------------- | ------------- | -------------------- |
| message         | module        | jinx                 |
| message         | sender        | `{sender address}`   |
| jinx_borrow     | borrow_coins  | `{amount}`           |
| jinx_withdrawal | borrower      | `{borrower address}` |

### MsgRepay

| Type       | Attribute Key | Attribute Value      |
| ---------- | ------------- | -------------------- |
| message    | module        | jinx                 |
| message    | sender        | `{sender address}`   |
| message    | owner         | `{owner address}`    |
| jinx_repay | repay_coins   | `{amount}`           |
| jinx_repay | sender        | `{borrower address}` |
