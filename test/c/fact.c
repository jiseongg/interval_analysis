int fact(int val) {
    int temp = 1;
    for (int i = 2; i <= val; i++)
        temp *= i;
    return temp;
}

int main() {
    fact(5);
}
