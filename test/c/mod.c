int mod (int a, int b) {
    int q = 0;
    int r = a;
    while (r >= b) {
        r = r - b;
        q = q + 1;
    }
    return r;
}
