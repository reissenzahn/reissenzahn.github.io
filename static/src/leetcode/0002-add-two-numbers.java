/**
 * Definition for singly-linked list.
 * public class ListNode {
 * int val;
 * ListNode next;
 * ListNode() {}
 * ListNode(int val) { this.val = val; }
 * ListNode(int val, ListNode next) { this.val = val; this.next = next; }
 * }
 */

// 11110
// 347
// 9685
// = 11032

// dummy -> 2 -> 3 -> 0 -> 1 -> 1
//                              ^ curr

class Solution {
  public ListNode addTwoNumbers(ListNode l1, ListNode l2) {
    final ListNode dummy = new ListNode();
    ListNode curr = dummy;

    int carry = 0;
    while (l1 != null || l2 != null || carry > 0) {  // note we cannot exit if carry == 1
      int sum = carry;

      if (l1 != null) {
        sum += l1.val;
        l1 = l1.next;
      }

      if (l2 != null) {
        sum += l2.val;
        l2 = l2.next;
      }

      // calculate next digit: 1X % 10 = X
      curr.next = new ListNode(sum % 10);
      curr = curr.next;

      // calculate carry: 1X / 10 = 1.X = 1
      carry = sum / 10;
    }

    return dummy.next;
  }
}

// Time complexity: O(max(n, m))
// Space complexity: O(1)