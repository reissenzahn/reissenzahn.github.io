class Solution {
  public int[] twoSum(int[] nums, int target) {
    // map elements of nums to their indicies
    final Map<Integer, Integer> m = new HashMap<>();

    for (int i = 0; i < nums.length; i++) {
      final int have = nums[i];

      // check if we have previously encountered a suitable complement
      final int want = target - have;

      if (m.containsKey(want)) {
        return new int[] { i, m.get(want) };
      }

      m.put(have, i);
    }

    // we are guaranteed that there will be a solution
    throw new IllegalArgumentException("Inputs do not satisfy target condition");
  }
}

// Time complexity: O(n)
// Space complexity: O(n)
