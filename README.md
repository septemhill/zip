# Repository Sample

This is a pretty simple repository design and implementation.

I created an interface called `Repository`, it would store different 
sub-repositories operations (e.g. `MutantRepository` and `OrderRepository`). 

In general, we could just create methods under `Resipository` interface 
(you could decide create another sub-repository to manage like sample or not). 
I also created an interface called `TransactionHolderRepository`, but for what ?

The thing is, for some scenario, we wouldn't have only one service. 
Services may use the sample package (e.g. `repository`) and each service has 
their own domain. If we have a serialized operations between services, and make 
sure each operation on specifed service is ready, then we would commit these
operation together, or rollback for all of them.

But `Repository` doesn't let caller own the transaction directly. 
So the `Repository` has `TransactionHolderRepository` which could create 
sub-repository with transaction, make caller could handle the transaction by their own.
